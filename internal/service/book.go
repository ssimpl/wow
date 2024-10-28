package service

import (
	"fmt"

	"github.com/ssimpl/wow/internal/domain"
)

type quoteRepo interface {
	GetQuoteByID(id int) (domain.Quote, error)
}

type Book struct {
	repo  quoteRepo
	curID int
}

func NewBook(repo quoteRepo) *Book {
	return &Book{
		repo: repo,
	}
}

func (b *Book) GetNextQuote() (domain.Quote, error) {
	quote, err := b.repo.GetQuoteByID(b.curID)
	if err != nil {
		return domain.Quote{}, fmt.Errorf("get quote from repo: %w", err)
	}

	b.curID++

	return quote, nil
}
