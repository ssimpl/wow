package handler

import (
	"encoding/json"
	"log/slog"
	"net"

	"github.com/ssimpl/wow/internal/domain"
)

const (
	messageTypeQuote     = "quote"
	messageTypeChallenge = "challenge"
	messageTypeError     = "error"
)

const bufferSize = 256

type powProvider interface {
	GenerateChallenge() (string, error)
	VerifyProof(challenge, proof string, difficulty int) bool
}

type quoteProvider interface {
	GetNextQuote() (domain.Quote, error)
}

type Handler struct {
	pow        powProvider
	quotes     quoteProvider
	difficulty int
}

func NewHandler(
	powProvider powProvider,
	quoteProvider quoteProvider,
	difficulty int,
) *Handler {
	return &Handler{
		pow:        powProvider,
		quotes:     quoteProvider,
		difficulty: difficulty,
	}
}

func (h *Handler) Handle(conn net.Conn) {
	defer conn.Close()

	slog.Info("Client has opened connection", "clientAddr", conn.RemoteAddr())

	challenge, err := h.pow.GenerateChallenge()
	if err != nil {
		errMsg := "Failed to generate challenge"
		slog.Error(errMsg, "error", err)
		sendError(conn, errMsg)
		return
	}

	slog.Info("Challenge was generated", "challenge", challenge)

	sendChallenge(conn, challenge, h.difficulty)

	buffer := make([]byte, bufferSize)
	n, err := conn.Read(buffer)
	if err != nil {
		errMsg := "Failed to read from client"
		slog.Error(errMsg, "error", err)
		sendError(conn, errMsg)
		return
	}

	proof := string(buffer[:n])

	slog.Info("Received a proof from client", "proof", proof)

	if h.pow.VerifyProof(challenge, proof, h.difficulty) {
		quote, err := h.quotes.GetNextQuote()
		if err != nil {
			errMsg := "Failed to get quote"
			slog.Error(errMsg, "error", err)
			sendError(conn, errMsg)
			return
		}
		sendQuote(conn, quote.Quote)
		slog.Info("Proof is verified and quote has been sent", "quote", quote)
	} else {
		errMsg := "Proof verification failed"
		slog.Warn(errMsg, "client", conn.RemoteAddr())
		sendError(conn, errMsg)
	}
}

func sendChallenge(conn net.Conn, challenge string, difficulty int) {
	sendResponse(conn, messageResponse{
		Type:       messageTypeChallenge,
		Message:    challenge,
		Difficulty: difficulty,
	})
}

func sendQuote(conn net.Conn, quote string) {
	sendResponse(conn, messageResponse{
		Type:    messageTypeQuote,
		Message: quote,
	})
}

func sendError(conn net.Conn, errStr string) {
	sendResponse(conn, messageResponse{
		Type:    messageTypeError,
		Message: errStr,
	})
}

func sendResponse(conn net.Conn, response messageResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		slog.Error("Failed to marshal message response", "error", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		slog.Error("Failed to write message response", "error", err)
		return
	}
}
