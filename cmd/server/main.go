package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"generate-leetcode/config"
	"generate-leetcode/internal/service"

	"github.com/sethvargo/go-envconfig"
)

func main() {
	var cfg config.Config
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatal("failed to process config:", err)
	}

	app, err := service.InitApp(ctx, &cfg)
	if err != nil {
		log.Fatal("failed to initialize app:", err)
	}

	go func() {
		slog.Info("Slack bot is running")
		if err := app.SocketClient.Run(); err != nil {
			log.Fatal("failed to run socket client:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Bot is shutting down...")
	cancel()
	slog.Info("Bot shut down successfully")
}
