package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ysuke25/hariko/internal/config"
	"github.com/ysuke25/hariko/internal/server"
)

var version = "dev"

func main() {
	configPath := flag.String("f", "routes.yaml", "path to config file")
	port := flag.Int("p", 0, "override server port")
	showVersion := flag.Bool("version", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Println("hariko version", version)
		os.Exit(0)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if *port > 0 {
		cfg.Server.Port = *port
	}

	srv := server.New(cfg, logger)
	if err := srv.Run(); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
