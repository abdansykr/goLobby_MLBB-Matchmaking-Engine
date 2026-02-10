package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/golobby/matchmaking/internal/config"
	"github.com/golobby/matchmaking/internal/database"
	"github.com/golobby/matchmaking/internal/delivery/http"
	"github.com/golobby/matchmaking/internal/repository"
	"github.com/golobby/matchmaking/internal/usecase"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("🚀 Starting GoLobby Matchmaking Engine...")
	log.Printf("Environment: %s", cfg.AppEnv)

	// Initialize PostgreSQL
	db, err := database.NewPostgresDB(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := database.NewRedisClient(cfg.GetRedisAddr(), cfg.Redis.Password)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize repositories
	teamRepo := repository.NewPostgresTeamRepository(db)
	matchRepo := repository.NewPostgresMatchRepository(db)
	cacheRepo := repository.NewRedisCache(redisClient)

	// Initialize usecase
	matchmakingConfig := usecase.MatchmakingConfig{
		InitialRankRange:  cfg.Matchmaking.InitialRankRange,
		ExtendedRankRange: cfg.Matchmaking.ExtendedRankRange,
		MatchTimeout:      cfg.Matchmaking.Timeout,
		ReadyTimeout:      cfg.Matchmaking.ReadyTimeout,
		GhostingPenalty:   cfg.Reputation.GhostingPenalty,
	}

	matchmakingUsecase := usecase.NewMatchmakingUsecase(
		teamRepo,
		matchRepo,
		cacheRepo,
		matchmakingConfig,
	)

	// Start matchmaking workers (4 workers by default)
	matchmakingUsecase.StartMatchmakingWorkers(4)
	defer matchmakingUsecase.StopMatchmakingWorkers()

	// Initialize Scrim repositories (NEW)
	scrimRequestRepo := repository.NewScrimRequestRepository(db)
	scrimMatchRepo := repository.NewScrimMatchRepository(db)
	rateLimiter := repository.NewRedisRateLimiter(redisClient)

	// Initialize Scrim usecase (NEW)
	scrimUsecase := usecase.NewScrimMatchmakingUsecase(
		scrimRequestRepo,
		scrimMatchRepo,
		rateLimiter,
		4, // 4 workers for scrim matching
	)

	// Start scrim matchmaking workers (NEW)
	scrimUsecase.Start()
	defer scrimUsecase.Stop()

	// Initialize WebSocket hub
	wsHub := http.NewWebSocketHub()

	// Start broadcast relay (from usecase to WebSocket clients)
	go func() {
		for msg := range matchmakingUsecase.GetBroadcastChannel() {
			// Relay to WebSocket hub
			hubMsg := &http.BroadcastMessage{
				Type:      msg.Type,
				TeamID:    msg.TeamID,
				MatchID:   msg.MatchID,
				Data:      msg.Data.(map[string]interface{}),
				Timestamp: msg.Timestamp.Format(time.RFC3339),
			}
			wsHub.Broadcast(hubMsg)
		}
	}()

	// Initialize handler
	handler := http.NewMatchmakingHandler(matchmakingUsecase, wsHub)

	// Initialize Scrim handler
	scrimHandler := http.NewScrimHandler(scrimUsecase, wsHub)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "GoLobby Matchmaking",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Health check
	app.Get("/health", handler.HealthCheck)

	// API routes
	api := app.Group("/api")

	// Matchmaking routes
	matchmaking := api.Group("/matchmaking")
	matchmaking.Post("/enqueue", handler.EnqueueTeam)
	matchmaking.Post("/ready", handler.ConfirmReady)
	matchmaking.Post("/cancel", handler.CancelMatchmaking)

	// Scrim routes
	http.RegisterScrimRoutes(app, scrimHandler)

	// WebSocket route
	app.Get("/ws", websocket.New(handler.WebSocketHandler))

	// Start server in a goroutine
	go func() {
		addr := ":" + cfg.AppPort
		log.Printf("🌐 Server starting on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("📴 Shutting down gracefully...")

	// Shutdown server
	if err := app.Shutdown(); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("👋 Server stopped")
}

