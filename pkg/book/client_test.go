package book

import (
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	srcChallenge := "12345"
	srcDifficulty := 4
	srcQuote := "The future belongs to those who believe in the beauty of their dreams."
	srcProof := "18083"

	t.Run("success", func(t *testing.T) {
		ln := startMockServer(t, srcChallenge, srcDifficulty, srcQuote, srcProof, false)

		client := NewClient(ln.Addr().String(), &powSolverMock{
			solveChallengeFunc: func(challenge string, difficulty int) string {
				assert.Equal(t, srcChallenge, challenge)
				assert.Equal(t, srcDifficulty, difficulty)
				return srcProof
			},
		}, 2*time.Second)
		resQuote, err := client.GetQuote()
		require.NoError(t, err)

		assert.Equal(t, srcQuote, resQuote)
	})
	t.Run("error", func(t *testing.T) {
		ln := startMockServer(t, srcChallenge, srcDifficulty, srcQuote, srcProof, true)

		client := NewClient(ln.Addr().String(), &powSolverMock{
			solveChallengeFunc: func(challenge string, difficulty int) string {
				assert.Equal(t, srcChallenge, challenge)
				assert.Equal(t, srcDifficulty, difficulty)
				return srcProof
			},
		}, 2*time.Second)
		_, err := client.GetQuote()
		require.Error(t, err)
		assert.ErrorContains(t, err, "Invalid proof")
	})
}

func startMockServer(
	t *testing.T, challenge string, difficulty int, quote string, srcProof string, serverError bool,
) net.Listener {
	ln, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go func() {
		conn, err := ln.Accept()
		require.NoError(t, err)
		defer conn.Close()

		challengeMsg := messageResponse{
			Type:       messageTypeChallenge,
			Message:    challenge,
			Difficulty: difficulty,
		}
		sendResponse(t, conn, challengeMsg)

		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)
		require.NoError(t, err)
		assert.Equal(t, srcProof, string(buffer[:n]))

		if serverError {
			errorMsg := messageResponse{
				Type:    messageTypeError,
				Message: "Invalid proof",
			}
			sendResponse(t, conn, errorMsg)
		} else {
			quoteMsg := messageResponse{
				Type:    messageTypeQuote,
				Message: quote,
			}
			sendResponse(t, conn, quoteMsg)
		}
	}()

	return ln
}

func sendResponse(t *testing.T, conn net.Conn, response messageResponse) {
	data, err := json.Marshal(response)
	require.NoError(t, err)

	_, err = conn.Write(data)
	require.NoError(t, err)
}

type powSolverMock struct {
	solveChallengeFunc func(challenge string, difficulty int) string
}

func (p *powSolverMock) SolveChallenge(challenge string, difficulty int) string {
	return p.solveChallengeFunc(challenge, difficulty)
}
