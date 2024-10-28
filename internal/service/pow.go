package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const challengeMaxNumber = 1_000_000

type POWProvider struct{}

func NewPOWProvider() *POWProvider {
	return &POWProvider{}
}

func (c *POWProvider) GenerateChallenge() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(challengeMaxNumber))
	if err != nil {
		return "", fmt.Errorf("get random number: %w", err)
	}
	return n.String(), nil
}

func (c *POWProvider) VerifyProof(challenge, proof string, difficulty int) bool {
	hash := sha256.Sum256([]byte(challenge + proof))
	hashStr := hex.EncodeToString(hash[:])
	return strings.HasPrefix(hashStr, strings.Repeat("0", difficulty))
}

func (c *POWProvider) SolveChallenge(challenge string, difficulty int) string {
	for nonce := 0; ; nonce++ {
		proof := strconv.Itoa(nonce)
		if c.VerifyProof(challenge, proof, difficulty) {
			return proof
		}
	}
}
