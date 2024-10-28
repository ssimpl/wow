package tcp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"
)

type handler interface {
	Handle(conn net.Conn)
}

type Server struct {
	addr                 string
	handler              handler
	clientWaitingTimeout time.Duration
	shutdownTimeout      time.Duration
	activeClients        sync.WaitGroup
}

func NewServer(addr string, handler handler, clientWaitingTimeout, shutdownTimeout time.Duration) *Server {
	return &Server{
		addr:                 addr,
		handler:              handler,
		clientWaitingTimeout: clientWaitingTimeout,
		shutdownTimeout:      shutdownTimeout,
	}
}

func (s *Server) Listen(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("start server: %w", err)
	}
	defer listener.Close()

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("Shutting down server")
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					slog.Error("Failed to accept connection", "error", err)
					continue
				}
				if err := conn.SetReadDeadline(time.Now().Add(s.clientWaitingTimeout)); err != nil {
					slog.Error("Failed to set read deadline", "error", err)
				}

				s.activeClients.Add(1)
				go s.handleConnection(conn)
			}
		}
	}()

	<-ctx.Done()

	stop := make(chan struct{})
	go func() {
		defer close(stop)
		s.activeClients.Wait()
	}()

	select {
	case <-stop:
	case <-time.After(s.shutdownTimeout):
	}

	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer s.activeClients.Done()
	s.handler.Handle(conn)
}
