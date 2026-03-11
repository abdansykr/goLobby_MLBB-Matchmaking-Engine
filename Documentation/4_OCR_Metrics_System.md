# OCR Validation & Metrics Architecture

## 1. OCR (Optical Character Recognition) System

Otomatisasi pengolahan tangkapan layar (Screenshot) hasil pertandingan Mobile Legends merupakan core spesifikasi yang diusung oleh Service tambahan (Microservice Python FastAPI + EasyOCR).

- **Fungsi Utama**: Mendeteksi teks "VICTORY" dan "DEFEAT" untuk verifikasi status penalti secara otonom.
- **Port Komunikasi**: `8001`
- **Tumpukan Teknologi**: GPU Accelerated OCR Scanner (`EasyOCR` Library), Numpy, FastAPI (Lightweight Python Router).

### Arsitektur Reputasi Match (Golang ke FastAPI)
```python
# FastAPI Endpoints
@app.post("/analyze-screenshot")
def analyze_screenshot(file: UploadFile):...
```
- Pemain yang membatalkan konfirmasi Match yang ditemukan pada **Match Found Modal** dikenakan *Penalty* point.
- Reputasi naik saat pertandingan tercatat (Screenshot sah tervalidasi).
- Reputasi turun jika terdapat *AFK / Pembatalan tiket*.

### Ekstraksi Teks

Algoritma pendeteksian FastAPI:
1. Skrip `main.py` menggunakan OpenCV merubah Citra (*Image*) Screenshot.
2. Filter kompresi mengisolasi area resolusi.
3. Neural Network merubah *Boundaries* menjadi raw teks.
4. "Defeat" dan "Victory" dikalkulasikan (Fuzzy Matching Logic) untuk memberi "Result Validation JSON" kembali ke Server Fiber Golang.

## 2. Prometheus & Grafana System

Untuk mengadopsi standar SLA produksi industri level-enterprise konvensional, GoLobby memperkenalkan *Observability* menggunakan stack Prom-Grafana:

### 1. Prometheus Scraper (`localhost:9090`)
File Konfigurasi YAML `monitoring/prometheus/prometheus.yml` di-mount ke kontainer Prometheus.

Metrik Go-Fiber disinari-keluar (di-*Expose*) via route HTTP `/metrics` dan secara sekuensial diekstrak setiap `5s`:

- **Active Matches Counter**: (`golobby_matchmaking_active_matches`) Metrik pengukur total permintaan matchmaking per detiknya.
- **Queue Length Gauge**: Pengukur total antrian yang mengendap atau menumpuk (bottleneck detection).

### 2. Grafana Dashboard (`localhost:3001`)

Memvisualisasi metrik *Real-Time*. Secara mandiri *Provisioning* dashboard GoLobby dibaca dari file `monitoring/grafana/dashboards/golobby.json` tanpa penyesuaian (*Out-of-the-Box Configuration*):

Tampilan (*Panels*):
1. Jumlah Pemain Mengantri (Real-time).
2. Sistem *Match Response Latency* API (Hitungan detak server).
3. POKE vs WARKOP Heatmap Analytics.
