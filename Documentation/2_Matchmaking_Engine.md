# Matchmaking Engine & Algorithms

## 1. Overview
Jantung aplikasi **GoLobby** adalah sub-sistem Matchmaking yang menggunakan struktur data Redis tingkat lanjut untuk antrian stateful. Ini memastikan tim yang dicari (POKE vs POKE) dipasangkan secara instan dan adil tanpa bottleneck *database locks* PostgreSQL SQL konvensional.

## 2. POKE vs WARKOP

Sistem mengklasifikasikan tiket berdasarkan 2 Mode (Category):

1. **POKE (Ranked Solo/Duo)**
   - Menggunakan mode *Skill-Based Matchmaking*.
   - Filter **Rank Tolerance**: ±1 Rank Weight.
   - Contoh: Jika ada pemain "Legend" mendaftar (Rank Weight = 6), sistem HANYA akan mencarikan musuh dengan Rank Weight 5 (Epic), 6 (Legend), atau 7 (Mythic).

2. **WARKOP (Pro Scrim)**
   - Menggunakan mode *Instant Join*.
   - WARKOP diatur pada Rank Weight khusus (Biasanya 10), sehingga menguasai prioritas antrian pro-team dan tidak memiliki filter rank ketat. Prioritasnya adalah kecepatan mencari musuh dengan sesama kategori `WARKOP`.

## 3. Worker Implementation (Go)

- **Redis Sorted Sets (ZSETS)**: Set ini dikomposisikan dari skor yang setara dengan `Waktu Entry` Tiket (`time.Now().UnixNano()`).
- Hal ini secara algoritmik memastikan antrian Matchmaking berjalan secara prinsip First-In-First-Out (FIFO) yang adil. Tiket terlama diantrikan di sisi terdepan.

### Alur Pencarian (*Loop Worker*)

```go
func (uc *matchmakingUsecase) StartMatchmakingWorker(ctx context.Context) {...}
```

1. **Pop Candidate Primer**: Worker terus meloop (Setiap detiknya mengambil minimal 1 antrian kosong teratas di Redis ZSET menggunakan fitur blokade data otomatis).
2. **Retrieve Range (Scanning)**: Setelah mendapatkan tiket Utama, Sistem melakukan *range query*. Range pencariannya di set dengan *Tolerance*. (Contoh: `Target Score = Weight ± Tolerance`).
3. **Locking Data (Redis Transactions / LUA Scripts)**: Ketika menemukan Kandidat 2 yang cocok, Worker langsung memberlakukan distributed *Pipelined Delete*/Transaksional. Ini mencegah **Race Condition**, dimana satu tim dicomot oleh dua match bertubrukan secara serentak.
4. **Data Persisted**: Match dibuat permanen `GenerateRandomMatchID()`. Data direkam ulang ke DB Historikal PostgreSQL.

## 4. Timeout System

- Jika sebuah tiket dibiarkan antri namun usianya melewati rentang `60 Detik` (Bisa dikonfigurasikan), tiket dihapus dari pool (*Garbage Collected*). 
- Notifikasi "Request Expired" kemudian dikirim ke Klien (via WebSocket atau Long-Polling).
