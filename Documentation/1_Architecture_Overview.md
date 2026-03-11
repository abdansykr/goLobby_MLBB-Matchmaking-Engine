# GoLobby - Architecture & System Overview

## 1. Executive Summary
**GoLobby** is a high-performance, real-time matchmaking engine designed for competitive mobile games (e.g., Mobile Legends). It pairs players and teams dynamically using rank-based tolerances and matchmaking algorithms backed by distributed queues.

## 2. Tech Stack Ecosystem
Sistem ini menggunakan arsitektur microservices-oriented dengan teknologi mutakhir:

- **Backend / Core Engine**: Golang 1.23 + Fiber Framework
- **Databases**: 
  - **PostgreSQL**: Relational database untuk Persistent Match History dan Teams.
  - **Redis**: In-Memory database untuk Stateful Matchmaking Queues, Caching, Pub/Sub, dan Rate Limiting.
- **Real-Time Communication**: Gorilla WebSockets (Go)
- **Frontend / UI**: Vue.js 3 + Vite + Vanilla Tailwind CSS (Custom Glassmorphism)
- **OCR Service (Validation)**: Python + FastAPI + EasyOCR
- **Monitoring & Metrics**: Prometheus & Grafana
- **Orchestration**: Docker & Docker Compose (Multi-stage builds)

## 3. High-Level Architecture Diagram

```mermaid
graph TD
    User([Player / Web UI - Vue.js])
    
    subgap [Nginx Reverse Proxy]
    end
    
    subgraph Core System [Backend - Go Fiber]
        API[REST API Handlers]
        WS[WebSocket Hub]
        Logic[Matchmaking Usecase]
        Repo[Data Repositories]
    end
    
    subgraph Data Layer
        PG[(PostgreSQL)]
        RD[(Redis)]
    end
    
    subgraph External Services
        OCR[Python OCR Service]
    end
    
    subgraph Monitoring
        Prom[Prometheus]
        Graf[Grafana]
    end
    
    User -->|HTTP /ws| WS
    User -->|HTTP POST| API
    
    API --> Logic
    WS --> Logic
    Logic --> Repo
    
    Repo -->|Persist History| PG
    Repo -->|Queue/Sorting/Cache| RD
    
    Logic -->|Submit Screenshot| OCR
    
    Prom -->|Scrape Metrics| API
    Graf -->|Query| Prom
```

## 4. Key Services Interaction
1. **Pendaftaran (Queueing)**: Frontend mengirim data tiket ke backend. Backend memilah tiket ke dalam Redis Sorted Set berdasarkan "Kategori" dan "Rank".
2. **Real-time Alert**: Begitu musuh ditemukan oleh *Matchmaking Worker*, backend menggunakan saluran WebSocket untuk mem-*push* notifikasi seketika ke kedua user.
3. **Penyimpanan Permanen**: Jika match terjadi (Matched), data disimpan dari Redis dialihkan permanen ke PostgreSQL.
4. **Validasi (OCR)**: Pemain dapat mengunggah bukti gambar pasca pertandingan. Backend akan meminta Python FastAPI untuk menganalisa gambar.
5. **Observability**: Prometheus terus-menerus mengambil metrik server (jumlah tiket aktif, metrik memory, latensi HTTP) menggunakan rute `/metrics`. Grafana menampilkan ini sebagai Dashboard visual.
