package handler

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/ssimpl/wow/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		powMock := &powProviderMock{
			challenge: "12345",
			verify:    true,
		}
		quoteMock := &quoteProviderMock{
			quote: domain.Quote{
				Quote:  "The future belongs to those who believe in the beauty of their dreams.",
				Author: "Eleanor Roosevelt",
			},
		}

		handler := NewHandler(powMock, quoteMock, 4)

		serverConn, clientConn := net.Pipe()
		defer serverConn.Close()
		defer clientConn.Close()

		go handler.Handle(serverConn)

		response := readResponse(t, clientConn)
		assert.Equal(t, messageResponse{
			Type:       messageTypeChallenge,
			Message:    powMock.challenge,
			Difficulty: 4,
		}, response)

		clientConn.Write([]byte("18083"))

		response = readResponse(t, clientConn)
		assert.Equal(t, messageResponse{
			Type:    messageTypeQuote,
			Message: quoteMock.quote.Quote,
		}, response)
	})
	t.Run("invalid proof", func(t *testing.T) {
		powMock := &powProviderMock{
			challenge: "12345",
			verify:    false,
		}
		quoteMock := &quoteProviderMock{
			quote: domain.Quote{
				Quote:  "The future belongs to those who believe in the beauty of their dreams.",
				Author: "Eleanor Roosevelt",
			},
		}

		handler := NewHandler(powMock, quoteMock, 4)

		serverConn, clientConn := net.Pipe()
		defer serverConn.Close()
		defer clientConn.Close()

		go handler.Handle(serverConn)

		response := readResponse(t, clientConn)
		assert.Equal(t, messageResponse{
			Type:       messageTypeChallenge,
			Message:    powMock.challenge,
			Difficulty: 4,
		}, response)

		clientConn.Write([]byte("1"))

		response = readResponse(t, clientConn)
		assert.Equal(t, messageResponse{
			Type:    messageTypeError,
			Message: "Proof verification failed",
		}, response)
	})
}

func readResponse(t *testing.T, conn net.Conn) messageResponse {
	buffer := make([]byte, 256)
	n, err := conn.Read(buffer)
	require.NoError(t, err)

	var response messageResponse
	err = json.Unmarshal(buffer[:n], &response)
	require.NoError(t, err)

	return response
}

type powProviderMock struct {
	challenge string
	verify    bool
	err       error
}

func (c *powProviderMock) GenerateChallenge() (string, error) {
	return c.challenge, c.err
}

func (c *powProviderMock) VerifyProof(challenge, proof string, difficulty int) bool {
	return c.verify
}

type quoteProviderMock struct {
	quote domain.Quote
	err   error
}

func (q *quoteProviderMock) GetNextQuote() (domain.Quote, error) {
	return q.quote, q.err
}
