# 🚀 GoLobby - MLBB Matchmaking Engine

A high-performance, specialized matchmaking system for Mobile Legends: Bang Bang scrims, built with Go.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Fiber](https://img.shields.io/badge/Fiber-v2-black?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=flat&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat&logo=docker)

## ✨ Statistics & Features

### **Matchmaking Logic**

The engine supports two distinct matchmaking categories with specialized rules:

#### 1. **POKE (Ranked Scrim)**
Designed for solo/duo players or casual teams looking for balanced matches.
*   **Ranks 1-8 (Warrior to Mythic Glory):**
    *   **Logic:** Strict `±1` rank tolerance.
    *   *Example:* A Rank 5 team will only match with Rank 4, 5, or 6.
*   **Rank 9 (Classic/Fun):**
    *   **Logic:** **Exact match only**.
    *   *Purpose:* Ensures casual players only match with other casual players for a fun experience.

#### 2. **WARKOP (Pro Scrim)**
Designed for competitive teams looking for instant matches regardless of rank.
*   **Logic:** **No rank tolerance**.
*   *Purpose:* Fast matching for pro teams who prioritize finding an opponent quickly over rank parity.

### **Core Systems**
*   **Auto-Scanner:** A background worker scans for matches every **2 seconds**, ensuring requests aren't left hanging.
*   **Request Expiry:** Scrim requests automatically expire after **60 minutes** to keep the pool fresh.
*   **Race Condition Safety:** Database-level concurrency control prevents double-booking of teams.
*   **REST API:** Clean, JSON-based API for frontend integration.

---

## 🏗️ Technical Architecture

### **Tech Stack**
*   **Language:** Go 1.21+
*   **Framework:** Fiber v2 (Fast HTTP Web Framework)
*   **Database:** PostgreSQL 15 (Data Persistence)
*   **Containerization:** Docker & Docker Compose
*   **Migration:** `golang-migrate`

### **Project Structure (Clean Architecture)**
```
matchMaking_go/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── domain/          # Business entities & interfaces
│   ├── repository/      # Database implementations (PostgreSQL)
│   ├── usecase/         # Matchmaking logic & business rules
│   ├── delivery/        # HTTP handlers
│   └── config/          # Environment configuration
├── migrations/          # Database schema migrations
└── frontend/            # Vue.js Frontend (Vite)
```

---

## 🚀 Quick Start

### **Prerequisites**
*   Docker & Docker Compose installed
*   Git

### **Run with Docker (Recommended)**
1.  **Clone the repository:**
    ```bash
    git clone https://github.com/abdansykr/goLobby_MLBB-Matchmaking-Engine.git
    cd matchMaking_go
    ```

2.  **Start the application:**
    ```bash
    docker-compose up --build -d
    ```

3.  **Access the application:**
    *   **Frontend:** `http://localhost:5173`
    *   **Backend API:** `http://localhost:3000`

### **Local Development (Manual)**
1.  **Database:** Start PostgreSQL.
2.  **Migrations:**
    ```bash
    migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/antigravity?sslmode=disable" up
    ```
3.  **Run Backend:**
    ```bash
    go run cmd/server/main.go
    ```

---

## 📡 API Endpoints

### **Scrim Requests**

#### **1. Create Scrim Request**
Enqueue a team for matchmaking.
- **URL:** `POST /api/scrim/request`
- **Body:**
  ```json
  {
    "team_name": "RRQ Hoshi",
    "whatsapp_number": "628123456789",
    "category": "POKE",  // or "WARKOP"
    "rank_weight": 5     // 1-10
  }
  ```

#### **2. Check Request Status**
Poll this endpoint to check if a match has been found.
- **URL:** `GET /api/scrim/request/:id`
- **Response (Matched):**
  ```json
  {
    "status": "matched",
    "match": {
      "opponent_team": "EVOS Legends",
      "whatsapp_url": "https://wa.me/628...",
      "match_id": "..."
    }
  }
  ```

#### **3. Cancel Request**
Remove a team from the matchmaking queue.
- **URL:** `POST /api/scrim/request/:id/cancel`

---

## 📝 Configuration (`.env`)

| Variable | Default | Description |
| :--- | :--- | :--- |
| `DB_HOST` | `postgres` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `postgres` | Database user |
| `DB_NAME` | `antigravity` | Database name |
| `SERVER_PORT` | `3000` | Backend API port |
| `ENABLE_RATE_LIMIT`| `true` | Enable API rate limiting |

---

## 🤝 Contributing

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
