package service

import (
	"testing"

	"github.com/ssimpl/wow/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBook_GetNextQuote(t *testing.T) {
	t.Run("get quote", func(t *testing.T) {
		repo := &quoteRepoMock{
			quote: domain.Quote{
				Quote:  "The future belongs to those who believe in the beauty of their dreams.",
				Author: "Eleanor Roosevelt",
			},
		}
		book := NewBook(repo)

		quote, err := book.GetNextQuote()
		require.NoError(t, err)

		assert.Equal(t, domain.Quote{
			Quote:  "The future belongs to those who believe in the beauty of their dreams.",
			Author: "Eleanor Roosevelt",
		}, quote)
	})
	t.Run("quote not found", func(t *testing.T) {
		repo := &quoteRepoMock{
			err: domain.ErrQuoteNotFound,
		}
		book := NewBook(repo)

		_, err := book.GetNextQuote()
		require.ErrorIs(t, err, domain.ErrQuoteNotFound)
	})
}

type quoteRepoMock struct {
	quote domain.Quote
	err   error
}

func (r *quoteRepoMock) GetQuoteByID(_ int) (domain.Quote, error) {
	return r.quote, r.err
}
