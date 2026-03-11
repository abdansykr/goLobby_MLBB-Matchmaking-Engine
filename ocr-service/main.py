"""
GoLobby OCR Microservice
FastAPI + Tesseract OCR for automated match result verification.

Flow:
1. Go backend uploads a screenshot via POST /verify
2. OCR extracts text from the image
3. Service parses winner, scores, and team names from the text
4. Confidence score is computed
5. Result is POSTed back to Go backend webhook
6. If confidence is low → match is flagged as "disputed"
"""

import asyncio
import io
import logging
import os
import re
import time
from enum import Enum
from typing import Optional

import cv2
import httpx
import numpy as np
import pytesseract
from fastapi import BackgroundTasks, FastAPI, File, Form, HTTPException, UploadFile
from fastapi.middleware.cors import CORSMiddleware
from PIL import Image
from prometheus_fastapi_instrumentator import Instrumentator
from pydantic import BaseModel
from pydantic_settings import BaseSettings

# ──────────────────────────────────────────────────────────────
# Configuration
# ──────────────────────────────────────────────────────────────

class Settings(BaseSettings):
    go_backend_url: str = "http://app:8080"
    go_webhook_secret: str = "ocr-webhook-secret-changeme"
    min_confidence: float = 60.0       # below → disputed
    tesseract_cmd: str = "/usr/bin/tesseract"
    log_level: str = "INFO"

    class Config:
        env_file = ".env"

settings = Settings()

# Configure Tesseract binary path
pytesseract.pytesseract.tesseract_cmd = settings.tesseract_cmd

# ──────────────────────────────────────────────────────────────
# Application
# ──────────────────────────────────────────────────────────────

logging.basicConfig(
    level=getattr(logging, settings.log_level),
    format='{"time":"%(asctime)s","level":"%(levelname)s","msg":"%(message)s"}',
)
logger = logging.getLogger("ocr-service")

app = FastAPI(
    title="GoLobby OCR Verification Service",
    description="AI-powered match result screenshot validator",
    version="1.0.0",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

# Attach Prometheus metrics endpoint at /metrics
Instrumentator().instrument(app).expose(app)


# ──────────────────────────────────────────────────────────────
# Pydantic Models
# ──────────────────────────────────────────────────────────────

class VerifyResult(str, Enum):
    VERIFIED = "verified"
    DISPUTED = "disputed"
    ERROR = "error"


class OCRResponse(BaseModel):
    match_id: str
    result: VerifyResult
    winner_team: Optional[str] = None
    score_team1: Optional[str] = None
    score_team2: Optional[str] = None
    confidence: float
    raw_text: str
    processing_ms: int


class WebhookPayload(BaseModel):
    match_id: str
    result: str
    winner_team: Optional[str]
    score_team1: Optional[str]
    score_team2: Optional[str]
    confidence: float
    raw_text: str


# ──────────────────────────────────────────────────────────────
# Image Pre-processing
# ──────────────────────────────────────────────────────────────

def preprocess_image(image_bytes: bytes) -> Image.Image:
    """
    Apply OpenCV pre-processing to improve OCR accuracy:
    - Convert to grayscale
    - Apply adaptive thresholding to handle varying lighting
    - Slight sharpening kernel
    """
    nparr = np.frombuffer(image_bytes, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

    # Grayscale
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

    # Upscale 2× for better OCR on small text
    gray = cv2.resize(gray, None, fx=2, fy=2, interpolation=cv2.INTER_CUBIC)

    # Adaptive threshold
    thresh = cv2.adaptiveThreshold(
        gray, 255,
        cv2.ADAPTIVE_THRESH_GAUSSIAN_C,
        cv2.THRESH_BINARY,
        11, 2
    )

    # Sharpening
    kernel = np.array([[0, -1, 0], [-1, 5, -1], [0, -1, 0]])
    sharpened = cv2.filter2D(thresh, -1, kernel)

    return Image.fromarray(sharpened)


# ──────────────────────────────────────────────────────────────
# Result Parsing
# ──────────────────────────────────────────────────────────────

MLBB_SCORE_PATTERN = re.compile(
    r"(\d{1,2})\s*[-:]\s*(\d{1,2})",  # matches "3-0", "3:0", "3 - 0"
)
WINNER_KEYWORDS = ["victory", "menang", "winner", "you won", "mvp"]
LOSER_KEYWORDS  = ["defeat", "kalah", "game over", "you lost"]


def parse_ocr_result(raw_text: str, team1_name: str, team2_name: str):
    """
    Heuristic parser for MLBB match result screenshots.
    Returns (winner_team, score_team1, score_team2, confidence).
    """
    text_lower = raw_text.lower()
    score_match = MLBB_SCORE_PATTERN.search(raw_text)

    score_team1 = score_match.group(1) if score_match else None
    score_team2 = score_match.group(2) if score_match else None

    winner_team = None
    confidence = 0.0

    # Check for known victory/defeat keywords
    has_victory = any(kw in text_lower for kw in WINNER_KEYWORDS)
    has_defeat  = any(kw in text_lower for kw in LOSER_KEYWORDS)

    # Check team names presence
    team1_found = team1_name.lower() in text_lower
    team2_found = team2_name.lower() in text_lower

    if score_match:
        confidence += 40.0  # found a score pattern

    if team1_found or team2_found:
        confidence += 20.0

    if has_victory or has_defeat:
        confidence += 30.0

    # Determine winner from score
    if score_match and score_team1 and score_team2:
        s1, s2 = int(score_team1), int(score_team2)
        if s1 > s2:
            winner_team = team1_name
            confidence += 10.0
        elif s2 > s1:
            winner_team = team2_name
            confidence += 10.0
        # Plausibility check: MLBB scrims are first to 3 wins
        if s1 > 5 or s2 > 5:
            confidence -= 20.0  # suspicious scores
            logger.warning("Suspicious score detected: %s-%s", score_team1, score_team2)

    return winner_team, score_team1, score_team2, min(confidence, 100.0)


# ──────────────────────────────────────────────────────────────
# Webhook Callback to Go Backend
# ──────────────────────────────────────────────────────────────

async def send_webhook(payload: WebhookPayload):
    """POST verification result back to the Go backend."""
    url = f"{settings.go_backend_url}/api/ocr/result"
    headers = {"X-OCR-Secret": settings.go_webhook_secret}
    async with httpx.AsyncClient(timeout=10) as client:
        try:
            resp = await client.post(url, json=payload.model_dump(), headers=headers)
            resp.raise_for_status()
            logger.info("Webhook sent successfully for match %s", payload.match_id)
        except httpx.HTTPError as exc:
            logger.error("Webhook failed for match %s: %s", payload.match_id, exc)


# ──────────────────────────────────────────────────────────────
# API Routes
# ──────────────────────────────────────────────────────────────

@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "golobby-ocr"}


@app.post("/verify", response_model=OCRResponse)
async def verify_screenshot(
    background_tasks: BackgroundTasks,
    match_id: str = Form(...),
    team1_name: str = Form(...),
    team2_name: str = Form(...),
    screenshot: UploadFile = File(...),
):
    """
    Upload a match result screenshot for AI verification.

    - **match_id**: UUID of the scrim_match record
    - **team1_name**: Name of team 1 (from the match record)
    - **team2_name**: Name of team 2
    - **screenshot**: Image file (PNG, JPG, WEBP accepted)
    """
    start_time = time.monotonic()

    # Validate file type
    allowed_types = {"image/png", "image/jpeg", "image/webp"}
    if screenshot.content_type not in allowed_types:
        raise HTTPException(
            status_code=400,
            detail=f"Unsupported file type: {screenshot.content_type}. Use PNG, JPEG, or WEBP.",
        )

    image_bytes = await screenshot.read()
    if len(image_bytes) > 10 * 1024 * 1024:  # 10 MB limit
        raise HTTPException(status_code=413, detail="Screenshot too large (max 10MB)")

    try:
        processed_img = preprocess_image(image_bytes)
        raw_text = pytesseract.image_to_string(processed_img, lang="eng+ind")
    except Exception as exc:
        logger.error("OCR failed for match %s: %s", match_id, exc)
        elapsed_ms = int((time.monotonic() - start_time) * 1000)
        return OCRResponse(
            match_id=match_id,
            result=VerifyResult.ERROR,
            confidence=0.0,
            raw_text="",
            processing_ms=elapsed_ms,
        )

    winner_team, score_team1, score_team2, confidence = parse_ocr_result(
        raw_text, team1_name, team2_name
    )

    result = VerifyResult.VERIFIED if confidence >= settings.min_confidence else VerifyResult.DISPUTED

    elapsed_ms = int((time.monotonic() - start_time) * 1000)
    logger.info(
        "OCR complete",
        extra={
            "match_id": match_id,
            "result": result,
            "confidence": confidence,
            "ms": elapsed_ms,
        },
    )

    payload = WebhookPayload(
        match_id=match_id,
        result=result.value,
        winner_team=winner_team,
        score_team1=score_team1,
        score_team2=score_team2,
        confidence=confidence,
        raw_text=raw_text[:2000],  # truncate for DB storage
    )

    # Fire-and-forget webhook – don't block the HTTP response
    background_tasks.add_task(send_webhook, payload)

    return OCRResponse(
        match_id=match_id,
        result=result,
        winner_team=winner_team,
        score_team1=score_team1,
        score_team2=score_team2,
        confidence=confidence,
        raw_text=raw_text[:2000],
        processing_ms=elapsed_ms,
    )
