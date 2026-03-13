package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/golobby/matchmaking/internal/config"
	"github.com/golobby/matchmaking/internal/database"
	httpHandler "github.com/golobby/matchmaking/internal/delivery/http"
	appLogger "github.com/golobby/matchmaking/internal/logger"
	"github.com/golobby/matchmaking/internal/repository"
	"github.com/golobby/matchmaking/internal/usecase"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func main() {
	// ── Configuration ──────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// ── Structured Logger ──────────────────────────────────────────────
	appLogger.Init(cfg.AppEnv)
	defer appLogger.Sync()

	appLogger.Info("🚀 Starting GoLobby Matchmaking Engine (Phase 2)",
		"env", cfg.AppEnv,
		"port", cfg.AppPort,
	)

	// ── PostgreSQL ─────────────────────────────────────────────────────
	db, err := database.NewPostgresDB(cfg.GetDatabaseDSN())
	if err != nil {
		appLogger.Fatal("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// ── Redis ──────────────────────────────────────────────────────────
	redisClient, err := database.NewRedisClient(cfg.GetRedisAddr(), cfg.Redis.Password)
	if err != nil {
		appLogger.Fatal("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	appLogger.Info("✅ Redis connected", "addr", cfg.GetRedisAddr())

	// ── Repositories ───────────────────────────────────────────────────
	teamRepo         := repository.NewPostgresTeamRepository(db)
	matchRepo        := repository.NewPostgresMatchRepository(db)
	cacheRepo        := repository.NewRedisCache(redisClient)
	scrimRequestRepo := repository.NewScrimRequestRepository(db)
	scrimMatchRepo   := repository.NewScrimMatchRepository(db)
	rateLimiter      := repository.NewRedisRateLimiter(redisClient)

	// ── Matchmaking Usecase ────────────────────────────────────────────
	matchmakingConfig := usecase.MatchmakingConfig{
		InitialRankRange:  cfg.Matchmaking.InitialRankRange,
		ExtendedRankRange: cfg.Matchmaking.ExtendedRankRange,
		MatchTimeout:      cfg.Matchmaking.Timeout,
		ReadyTimeout:      cfg.Matchmaking.ReadyTimeout,
		GhostingPenalty:   cfg.Reputation.GhostingPenalty,
	}
	matchmakingUsecase := usecase.NewMatchmakingUsecase(
		teamRepo, matchRepo, cacheRepo, matchmakingConfig,
	)
	matchmakingUsecase.StartMatchmakingWorkers(4)
	defer matchmakingUsecase.StopMatchmakingWorkers()

	// ── Scrim Usecase ──────────────────────────────────────────────────
	scrimUsecase := usecase.NewScrimMatchmakingUsecase(
		scrimRequestRepo, scrimMatchRepo, rateLimiter, cacheRepo, 4,
	)
	scrimUsecase.Start()
	defer scrimUsecase.Stop()

	// ── WebSocket Hub ──────────────────────────────────────────────────
	wsHub := httpHandler.NewWebSocketHub()

	// Relay matchmaking events → WebSocket clients
	go func() {
		for msg := range matchmakingUsecase.GetBroadcastChannel() {
			hubMsg := &httpHandler.BroadcastMessage{
				Type:      msg.Type,
				TeamID:    msg.TeamID,
				MatchID:   msg.MatchID,
				Data:      msg.Data.(map[string]interface{}),
				Timestamp: msg.Timestamp.Format(time.RFC3339),
			}
			wsHub.Broadcast(hubMsg)
		}
	}()

	// ── Handlers ───────────────────────────────────────────────────────
	matchmakingHandler := httpHandler.NewMatchmakingHandler(matchmakingUsecase, wsHub)
	scrimHandler       := httpHandler.NewScrimHandler(scrimUsecase, wsHub)
	ocrHandler         := httpHandler.NewOCRHandler(wsHub)

	// ── Fiber App ──────────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName:      "GoLobby Matchmaking v2",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: `{"time":"${time}","status":${status},"latency":"${latency}","method":"${method}","path":"${path}"}` + "\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-OCR-Secret",
	}))

	// ── Routes ─────────────────────────────────────────────────────────
	app.Get("/health", matchmakingHandler.HealthCheck)

	api := app.Group("/api")

	matchmaking := api.Group("/matchmaking")
	matchmaking.Post("/enqueue", matchmakingHandler.EnqueueTeam)
	matchmaking.Post("/ready",   matchmakingHandler.ConfirmReady)
	matchmaking.Post("/cancel",  matchmakingHandler.CancelMatchmaking)

	httpHandler.RegisterScrimRoutes(app, scrimHandler)
	httpHandler.RegisterOCRRoutes(app, ocrHandler)

	// WebSocket upgrade
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws", websocket.New(matchmakingHandler.WebSocketHandler))

	// ── Prometheus Metrics Server (separate port) ──────────────────────
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "2112"
	}
	go func() {
		metricsApp := fiber.New(fiber.Config{DisableStartupMessage: true})
		metricsApp.Get("/metrics", func(c *fiber.Ctx) error {
			fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())(c.Context())
			return nil
		})
		appLogger.Info("📊 Prometheus metrics listening", "port", metricsPort)
		if err := metricsApp.Listen(":" + metricsPort); err != nil {
			appLogger.Error("Metrics server error", "err", err)
		}
	}()

	// ── Start Main Server ──────────────────────────────────────────────
	go func() {
		addr := ":" + cfg.AppPort
		appLogger.Info("🌐 HTTP server starting", "addr", addr)
		if err := app.Listen(addr); err != nil {
			appLogger.Fatal("Failed to start server: %v", err)
		}
	}()

	// ── Graceful Shutdown ──────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("📴 Shutting down gracefully...")
	if err := app.Shutdown(); err != nil {
		appLogger.Error("Server shutdown error", "err", err)
	}
	appLogger.Info("👋 Server stopped")
}
