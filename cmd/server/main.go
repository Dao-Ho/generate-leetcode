package main

import (
	"context"
	"log"

	"generate-leetcode/config"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	var cfg config.Config
	ctx := context.Background()
	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatal("failed to process config:", err)
	}
}
