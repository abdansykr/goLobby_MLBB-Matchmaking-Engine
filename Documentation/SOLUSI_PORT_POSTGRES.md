# 🔧 Solusi Konflik Port PostgreSQL

## ❌ Masalah

Anda memiliki PostgreSQL 17 yang terinstall di local dan sudah menggunakan port **5432**, sehingga terjadi konflik ketika Docker container GoLobby mencoba menggunakan port yang sama.

```
Error: Bind for 0.0.0.0:5432 failed: port is already allocated
```

---

## ✅ Solusi yang Telah Diterapkan

PostgreSQL container GoLobby telah dikonfigurasi untuk menggunakan **port 5433** (bukan 5432) agar tidak konflik dengan PostgreSQL 17 local Anda.

### Perubahan yang Dilakukan:

1. **docker-compose.yaml**
   ```yaml
   ports:
     - "5433:5432"  # Host port 5433 → Container port 5432
   ```

2. **.env**
   ```bash
   DB_PORT=5433
   ```

---

## 📊 Mapping Port Terkini

| Service | Local Port | Container Port | Status |
|---------|------------|----------------|--------|
| **PostgreSQL Local** | 5432 | - | PostgreSQL 17 (Anda) |
| **GoLobby Postgres** | **5433** | 5432 | postgres:15-alpine (Docker) |
| **GoLobby Redis** | 6379 | 6379 | redis:7-alpine (Docker) |
| **GoLobby App** | 3000 | 8080 | Go Application (Docker) |

---

## 🔌 Cara Koneksi ke Database

### 1. Koneksi ke PostgreSQL GoLobby (Container)

**Via Docker Exec:**
```powershell
docker exec -it golobby_postgres psql -U postgres -d golobby
```

**Via psql dari Local (jika psql terinstall):**
```powershell
psql -h localhost -p 5433 -U postgres -d golobby
```

**Connection String:**
```
postgresql://postgres:postgres@localhost:5433/golobby
```

**PgAdmin / DBeaver:**
- Host: `localhost`
- Port: `5433`
- Username: `postgres`
- Password: `postgres`
- Database: `golobby`

---

### 2. Koneksi ke PostgreSQL 17 Local (Anda)

Tetap menggunakan port **5432** seperti biasa:

```powershell
psql -h localhost -p 5432 -U postgres -d your_database
```

**Connection String:**
```
postgresql://postgres:password@localhost:5432/your_database
```

---

## 🧪 Testing Kedua Database

### Test PostgreSQL GoLobby (Port 5433)

```powershell
# Cek versi
docker exec golobby_postgres psql -U postgres -c "SELECT version();"

# Lihat tables
docker exec golobby_postgres psql -U postgres -d golobby -c "\dt"

# Lihat data teams
docker exec golobby_postgres psql -U postgres -d golobby -c "SELECT * FROM teams;"
```

### Test PostgreSQL 17 Local (Port 5432)

```powershell
# Cek versi (harus PostgreSQL 17)
psql -h localhost -p 5432 -U postgres -c "SELECT version();"
```

---

## 🔄 Alternatif Solusi Lain

### Opsi A: Stop PostgreSQL Local Service (Tidak Recommended)

Jika Anda tidak menggunakan PostgreSQL 17 local, bisa stop servicenya:

```powershell
# Stop PostgreSQL service
Stop-Service postgresql-x64-17

# Kembalikan port Docker ke 5432
# Edit docker-compose.yaml: - "5432:5432"
```

### Opsi B: Upgrade Docker Image ke PostgreSQL 17

Edit `docker-compose.yaml`:
```yaml
postgres:
  image: postgres:17-alpine  # Dari 15-alpine ke 17-alpine
```

Lalu rebuild:
```powershell
docker-compose down -v
docker-compose up -d --build
```

⚠️ **Warning**: Ini akan hapus semua data di container!

---

## 📝 Catatan Penting

1. **Data Terpisah**: Data di PostgreSQL GoLobby (container) terpisah dari PostgreSQL 17 local Anda

2. **Persistence**: Data PostgreSQL GoLobby disimpan di Docker volume `matchmaking_go_postgres_data`

3. **Backup Data**:
   ```powershell
   # Backup GoLobby database
   docker exec golobby_postgres pg_dump -U postgres golobby > backup_golobby.sql
   
   # Restore
   docker exec -i golobby_postgres psql -U postgres golobby < backup_golobby.sql
   ```

4. **Reset Database**:
   ```powershell
   # Hapus semua data dan start fresh
   docker-compose down -v
   docker-compose up -d
   ```

---

## ✅ Verifikasi Setup

Cek semua port listening:

```powershell
# PostgreSQL ports
netstat -ano | findstr :5432  # PostgreSQL 17 Local
netstat -ano | findstr :5433  # PostgreSQL GoLobby

# Application
netstat -ano | findstr :3000  # GoLobby App

# Redis
netstat -ano | findstr :6379  # GoLobby Redis
```

---

## 🚀 Quick Commands

```powershell
# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker logs golobby_postgres
docker logs golobby_app

# Connect to database
docker exec -it golobby_postgres psql -U postgres -d golobby

# Stop all
docker-compose down
```

---

## 🎯 Status Terkini

✅ **PostgreSQL 17 Local**: Running di port **5432**  
✅ **PostgreSQL GoLobby**: Running di port **5433**  
✅ **GoLobby App**: Running di port **3000**  
✅ **Redis**: Running di port **6379**  

**Tidak ada konflik port! Semua berjalan sempurna! 🎉**
