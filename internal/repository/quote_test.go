package repository

import (
	"testing"

	"github.com/ssimpl/wow/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuote_GetQuoteByID(t *testing.T) {
	t.Run("get quote", func(t *testing.T) {
		repo, err := NewQuote()
		require.NoError(t, err)

		quote, err := repo.GetQuoteByID(3)
		require.NoError(t, err)

		assert.Equal(t, domain.Quote{
			Quote:  "The future belongs to those who believe in the beauty of their dreams.",
			Author: "Eleanor Roosevelt",
		}, quote)
	})
	t.Run("id out of range", func(t *testing.T) {
		repo, err := NewQuote()
		require.NoError(t, err)

		quote, err := repo.GetQuoteByID(29)
		require.NoError(t, err)

		assert.Equal(t, domain.Quote{
			Quote:  "The way to get started is to quit talking and begin doing.",
			Author: "Walt Disney",
		}, quote)
	})
	t.Run("quote not found", func(t *testing.T) {
		sourceQuotes = []byte(`[]`)

		repo, err := NewQuote()
		require.NoError(t, err)

		_, err = repo.GetQuoteByID(1)
		require.ErrorIs(t, err, domain.ErrQuoteNotFound)
	})
}
