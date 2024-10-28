package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ssimpl/wow/internal/repository"
	"github.com/ssimpl/wow/internal/service"
	"github.com/ssimpl/wow/internal/transport/tcp"
	"github.com/ssimpl/wow/internal/transport/tcp/handler"
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

	quoteRepo, err := repository.NewQuote()
	if err != nil {
		return fmt.Errorf("create quote repository: %w", err)
	}

	powProvider := service.NewPOWProvider()
	book := service.NewBook(quoteRepo)
	handler := handler.NewHandler(powProvider, book, cfg.Difficulty)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	server := tcp.NewServer(cfg.Addr, handler, cfg.ClientWaitingTimeout, cfg.ShutdownTimeout)
	if err := server.Listen(ctx); err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}
