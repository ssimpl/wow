package book

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"time"
)

const bufferSize = 256

type powSolver interface {
	SolveChallenge(challenge string, difficulty int) string
}

type Client struct {
	serverAddr string
	powSolver  powSolver
	timeout    time.Duration
}

func NewClient(serverAddr string, powSolver powSolver, timeout time.Duration) *Client {
	return &Client{
		serverAddr: serverAddr,
		powSolver:  powSolver,
		timeout:    timeout,
	}
}

func (c *Client) GetQuote() (string, error) {
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		return "", fmt.Errorf("connect to server: %w", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
		return "", fmt.Errorf("set deadline: %w", err)
	}

	slog.Info("Client is connected to server", "serverAddr", conn.RemoteAddr())

	buffer := make([]byte, bufferSize)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("read server response: %w", err)
	}

	response, err := parseResponse(buffer[:n])
	if err != nil {
		return "", fmt.Errorf("parse server response: %w", err)
	}

	if response.Type != messageTypeChallenge {
		return "", fmt.Errorf("unexpected server response type: %q", response.Type)
	}

	slog.Info("Challenge was received from server", "challenge", response.Message, "difficulty", response.Difficulty)

	proof := c.powSolver.SolveChallenge(response.Message, response.Difficulty)

	slog.Info("Challenge is solved", "proof", proof)

	if _, err := conn.Write([]byte(proof)); err != nil {
		return "", fmt.Errorf("send proof: %w", err)
	}

	n, err = conn.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("read server response: %w", err)
	}

	response, err = parseResponse(buffer[:n])
	if err != nil {
		return "", fmt.Errorf("parse server response: %w", err)
	}

	switch response.Type {
	case messageTypeQuote:
		return response.Message, nil
	case messageTypeError:
		return "", fmt.Errorf("server error: %s", response.Message)
	default:
		return "", fmt.Errorf("unexpected server response type: '%q'", response.Type)
	}
}

func parseResponse(data []byte) (messageResponse, error) {
	var response messageResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return response, fmt.Errorf("unmarshal server response: %w", err)
	}
	return response, nil
}
