package port

import (
	"context"
	"personal-finance/core/domain"
)

type TransactionRepository interface {
	GetTransactions(ctx context.Context, page, limit uint64) ([]domain.Transaction, any, any, error)
	GetTransaction(ctx context.Context, id string) (*domain.Transaction, error)
	CreateTransaction(ctx context.Context, createTransaction *domain.Transaction) (*domain.Transaction, error)
	UpdateTransaction(ctx context.Context, id string, updatedTransaction *domain.Transaction) (*domain.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
}

type TransactionService interface {
	GetStatus(ctx context.Context) string
	GetTransactions(ctx context.Context, page, limit uint64) ([]domain.Transaction, any, any, error)
	GetTransaction(ctx context.Context, id string) (*domain.Transaction, error)
	CreateTransaction(ctx context.Context, createTransaction *domain.Transaction) (*domain.Transaction, error)
	UpdateTransaction(ctx context.Context, id string, updatedTransaction *domain.Transaction) (*domain.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
}
