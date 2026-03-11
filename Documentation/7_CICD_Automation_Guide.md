# GoLobby - Panduan CI/CD & Minimal Downtime 🚀

Dokumen ini berisi informasi operasional untuk mempersiapkan otomatisasi penerapan (Deployment) `GitHub Actions` di server VPS.

## 1. Persiapan Awal Server (VPS)

Sebelum *GitHub Actions Pipeline* bisa berjalan secara independen, Server VPS production mu harus disiapkan sekali untuk selamanya. Gunakan perintah SSH ini:

### Langkah A. Mendaftar SSH Key
Buat *Private & Public Key* di VPS-mu untuk digunakan oleh Action. Data Privatnya nanti disalin ke Settings -> Secrets pada Repository Git.

```bash
# Masuk ke server VPS utama mu
ssh user@IP_SERVER

# Jalankan perintah ini (Tekan Enter saja bila diminta lokasi dan password kosong)
ssh-keygen -t rsa -b 4096 -C "deploy@golobby"

# Lihat dan Copy seluruh isi Private Key (Dimulai --BEGIN RSA PRIVATE KEY--)
cat ~/.ssh/id_rsa
```

### Langkah B. Inisiasi Folder Repositori Awal

Buat folder master tempat kode aplikasi dijalankan, di mana GitHub akan melabuhkan (pull) pembaruan di masa depan.
Lalu install Docker + Docker Compose.

```bash
# Pastikan Git dan Docker sudah terpasang
sudo apt update && sudo apt install git docker.io docker-compose-v2 -y

# Buat direktori /var/www atau sesuaikan dengan keinginan kalian
mkdir -p /var/www

# Clone repositori awal dari github
cd /var/www/
git clone https://github.com/abdansykr/goLobby_MLBB-Matchmaking-Engine.git goLobby

# Buat file konfigurasi rahasia pertama kalinya 
cd goLobby
# Konfigurasikan file env
cp .env.example .env
nano .env # (Isi sesuai production server, e.g. OCR URL dan Webhook)
```

## 2. Cara Mengatasi "Downtime Minimal" (Zero/Minimal Downtime Deployment)

Saat `docker compose up -d --build` berjalan pada tahap CD (Continuous Deployment), aplikasi bisa saja mengalami mati sebentar ('Downtime') saat versi baru mematikan yang lama.

Inilah mengapa di file `deploy.yml`, tahapan CI/CD dipisah:

1. **Tahap Verifikasi**: GitHub me-lint dan mengetes kode Golang (`go test ./...`) dan mensimulasikan `docker build` pada server runner GitHub (Bukan VPS kamu). Jika rusak/salah syntax, Build dibatalkan dan mesin VPS yang sedang online tetap melayani player *Matchmaking* GoLobby tanpa terganggu.
2. **Tahap Detached Mode (`-d`)**: Skrip `docker compose up -d --build` memiliki fitur cerdas (*Smart Recreation*). Docker tidak mematikan service Go atau Postgres terlebih dahulu secara serentak. Ia hanya mem-*build image* barunya di belakang layar. Kontainer versi lama secara harfiah akan dihancurkan bersamaan di detik kontainer *versi baru* selesai tercipta.

**Pengaturan Khusus Zero-Downtime Tanpa Terputus Sama Sekali:**
Jika trafik *scrim/mabar* GoLobby sungguh padat dan pembaruan kontainer tidak boleh terputus barang 1 milidetik pun, gunakan opsi Load Balancer dan Health-Checks *(Blue-Green Deployment)* di masa mendatang:

```yaml
# Tambahan di docker-compose.yml 
services:
  app:
    # Memerintahkan Docker untuk memutar instansi baru SEBELUM menyapu yang lama
    deploy:
      update_config:
        order: start-first
```

## 3. Tambahan Keamanan
Pastikan 3 Variabel Utama dimasukkan pada Menu GitHub:
`Settings > Secrets and variables > Actions > New repository secret`:
1. `SERVER_IP` (Ex: `103.18.xx.xx`)
2. `SERVER_USER` (Ex: `root` atau `ubuntu`)
3. `SSH_PRIVATE_KEY` (Isi *Private RSA Key* dari Langkah 1A di atas).
