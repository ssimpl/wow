package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChallenge(t *testing.T) {
	t.Run("generate challenge", func(t *testing.T) {
		provider := NewPOWProvider()

		challenge, err := provider.GenerateChallenge()
		require.NoError(t, err)
		assert.NotEmpty(t, challenge)

		challenge2, err := provider.GenerateChallenge()
		require.NoError(t, err)
		assert.NotEmpty(t, challenge)

		assert.NotEqual(t, challenge, challenge2)
	})
	t.Run("success proof", func(t *testing.T) {
		provider := NewPOWProvider()

		challenge := "12345"
		proof := "18083"

		assert.True(t, provider.VerifyProof(challenge, proof, 4))
	})
	t.Run("invalid proof", func(t *testing.T) {
		provider := NewPOWProvider()

		challenge := "12345"
		proof := "1"

		assert.False(t, provider.VerifyProof(challenge, proof, 4))
	})
	t.Run("solve challenge", func(t *testing.T) {
		provider := NewPOWProvider()

		challenge := "12345"
		difficulty := 4

		proof := provider.SolveChallenge(challenge, difficulty)
		assert.Equal(t, "18083", proof)
	})
}
