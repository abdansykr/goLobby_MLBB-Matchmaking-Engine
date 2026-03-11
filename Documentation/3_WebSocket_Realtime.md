# WebSocket & Real-Time Notification System

## 1. Desain Jaringan
Aplikasi ini berevolusi dari metode konvensional *Long-Polling* menuju jalur bi-directional WebSockets (`ws://`) untuk menjamin kecepatan komunikasi data latensi ultra-rendah setelah lawan ditentukan oleh Redis Matchmaking Worker.

Komunikasi Websocket beroperasi secara simultan di Fiber dengan package Go `github.com/gofiber/websocket/v2`.

## 2. Websocket Hub
Server mengatur koneksi berjalan via "Connection Hub" yang terpusat di `scrim_handler.go`:
```go
type ScrimHandler struct {
    sync.RWMutex
    Connections map[string]*websocket.Conn
}
```

Dalam hub ini, Thread Safety dikontrol melalui struktur *RWMutex*. Map (`Connections`) berisi `request_id` yang me-mapping setiap channel ws aktif dengan klien spesifik. Pemain bisa secara asinkron menunggu "Sinyal Server".

## 3. Alur Komunikasi

1. **Client Start**: User menekan *Find Match*.
2. **API Request**: HTTP POST (`/api/scrim/request`) menyimpan data, lalu me-return sebuah `request_id`.
3. **Upgrade WS Connection**: Klien web (Vue.js) kemudian merubah jalur (`ws://localhost/ws?request_id=xxx`), server Fiber men-*Upgrade* koneksi tersebut menjadi Full Duplex.
4. **Listener Asinkron**: Thread *Goroutine Connection* didaftarkan ke map. Sisi klien mendengarkan `onmessage()`.
5. **Worker Execution**: Subsistem Redis *Matchmaking Engine Usecase* menemukan kecocokan antar Tiket "A" dan Tiket "B".
6. **Broadcasting Saluran Aktif**: Worker mengeksekusi metode injeksi pada REST layer:
   ```go
   h.NotifyMatchFound(matchResult) 
   ```
7. **Pesan JSON WS Diterima**: WS Backend langsung mengirimkan JSON Serialized `type: SCRIM_MATCH_FOUND` kepada Klien A dan Klien B.
8. **Teardown**: Setelah sukses disiarkan, Websocket yang aktif secara graceful dihentikan `Disconnect()`.

## 4. Keamanan dan Fallback

- Hub akan membersihkan *Zombie Connections* (contoh: user menutup browser secara fisik).
- Klien dirancang tangguh (*Resilient*): Jika browser user memblokir port Websocket / koneksi internet drop, Frontend memiliki metode cadangan (*Polling Graceful Degradation*) yang mundur perlahan mengirimkan GET HTTP request per-3 detik hingga menemukan tiket di PostgreSQL.
