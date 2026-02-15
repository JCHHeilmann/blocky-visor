package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/JCHHeilmann/blocky-visor/sidecar/handler"
	"github.com/JCHHeilmann/blocky-visor/sidecar/middleware"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS(cfg.CORSOrigins))

	// Health check — no auth
	r.Get("/api/health", handler.Health)

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.APIKeyAuth(cfg.APIKey))

		r.Get("/api/config", handler.GetConfig(cfg.Blocky.ConfigPath))
		r.Put("/api/config", handler.PutConfig(cfg.Blocky.ConfigPath))

		r.Get("/api/service/status", handler.ServiceStatus(cfg.Blocky.ServiceName))
		r.Post("/api/service/restart", handler.ServiceRestart(cfg.Blocky.ServiceName))

		r.Get("/api/stats", handler.GetStats(cfg.Blocky.LogDir))
		r.Get("/api/stats/timeline", handler.GetTimeline(cfg.Blocky.LogDir))

		r.Get("/api/logs", handler.GetLogs(cfg.Blocky.LogDir))
	})

	fmt.Printf("Blocky Visor sidecar listening on %s\n", cfg.Listen)
	if err := http.ListenAndServe(cfg.Listen, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
