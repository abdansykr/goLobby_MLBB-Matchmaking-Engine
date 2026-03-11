# Containerization & Deployment

## 1. Arsitektur Docker Compose
Proyek **GoLobby** seluruhnya ditenagai oleh Docker, mengabstraksi manajemen instalasi dependen lokal OS (Go, Node.js, PostgreSQL). Docker Compose secara orkestrasi menyatukan 7 container yang berdiri masing-masing ke dalam satu "Virtual Network" yang saling terkait.

## 2. Layanan (Services) / Komponen `docker-compose.yaml`

1. **`app`** (Backend API - Golang Fiber)
   - Port Publik: `3000`
   - File Dockerfile: `. /Dockerfile` berbasis `golang:1.23-alpine`. 
   - Di-*build* multi stage memisahkan runtime dengan compiler Go.

2. **`frontend`** (UI - Vue.js & Nginx)
   - Port Publik: `5173`
   - Direktori `frontend/Dockerfile` Multi-stage mengkompilasi file statis Vite (`npm run build`) via image node dan melayani `dist/` index secara asinkron menggunakan base OS Alpine **Nginx**.
   - Berfungsi melayani Static HTTP requests & proxy layer.

3. **`ocr-service`** (Python Validation API)
   - Port Publik: `8001`
   - Image FastAPI berbasis Python. Keperluan *Machine Learning EasyOCR* dan pengolahan *Computer Vision OpenCV*. (Kinerja image ini sengaja diisolasi supaya tidak memakan thread RAM Go app).

4. **`postgres`** (Database Relasional Transaksi)
   - Port Publik: `5433` (Dipetakan ulang dari standar `5432` agar tak konflik).
   - Di-mount oleh initial schema SQL `migrations/` & persistent Volume Docker untuk menahan data user saat docker dimatikan.

5. **`redis`** (In-Memory Data Store)
   - Port Lokal: `6379`.
   - Mengontrol algoritma antrian Matchmaking real-time dan cache. `alpine` image base untuk performa ringan.

6. **`prometheus`** & 7. **`grafana`** (Monitoring Suite)
   - Port: `9090` (Prom) dan `3001` (Grafana). Menyedot log dan data uptime sistem.

## 3. Proses Build Terotomatisasi (Perbaikan Sistem)
Peningkatan signifikan telah dilakukan pada Dockerfile:

### 1. Perbaikan Ketergantungan Versi `go 1.25.5`
- Sebelumnya sistem Golang Windows mengalami `corrupt index`. Solusi dalam `Dockerfile` di-lock ke stabil: `FROM golang:1.23-alpine`.
- Perintah kompilasi utama diperbarui ke standar produksi linux stabil.

### 2. Isu Konflik Port Jaringan
- Sebelumnya Vue Frontend (3000) dan Backend (3000) tabrakan di OS lokal.
- Solusi: `docker-compose.yaml` dialihkan `frontend: "5173:80"`. Menyelesaikan isu operasional pada env localhost Windows. Klien dan API sekarang berjalan secara kooperatif, diatur proxy Nginx (80) pada bridge network internal docker Compose `golobby-network`.
