package repository

import (
	_ "embed"
	"encoding/json"

	"github.com/ssimpl/wow/internal/domain"
)

//go:embed quotes.json
var sourceQuotes []byte

type Quote struct {
	quotes []quoteEntity
}

func NewQuote() (*Quote, error) {
	var quotes []quoteEntity
	if err := json.Unmarshal(sourceQuotes, &quotes); err != nil {
		return nil, err
	}
	return &Quote{quotes: quotes}, nil
}

func (q *Quote) GetQuoteByID(id int) (domain.Quote, error) {
	if len(q.quotes) == 0 {
		return domain.Quote{}, domain.ErrQuoteNotFound
	}

	i := id % len(q.quotes)
	return domain.Quote{
		Quote:  q.quotes[i].Quote,
		Author: q.quotes[i].Author,
	}, nil
}
