package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/ssimpl/wow/internal/service"
	"github.com/ssimpl/wow/pkg/book"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := newConfig()
	if err != nil {
		return fmt.Errorf("create config: %w", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting a client", "serverAddr", cfg.ServerAddr, "timeout", cfg.ConnectionTimeout)

	powProvider := service.NewPOWProvider()
	client := book.NewClient(cfg.ServerAddr, powProvider, cfg.ConnectionTimeout)

	wg := sync.WaitGroup{}
	for i := 0; i < cfg.Requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			quote, err := client.GetQuote()
			if err != nil {
				slog.Error("Failed to get quote", "error", err)
				return
			}

			slog.Info("Quote was received from server", "quote", quote)
		}()

	}
	wg.Wait()

	return nil
}
