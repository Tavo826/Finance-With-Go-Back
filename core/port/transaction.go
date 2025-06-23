package port

import (
	"context"
	"personal-finance/core/domain"
)

type TransactionRepository interface {
	GetTransactionsByUserId(ctx context.Context, page, limit uint64, userId string) ([]domain.Transaction, any, any, error)
	GetTransactionsByDate(ctx context.Context, userId string, page, limit uint64, year int, month int) ([]domain.Transaction, any, any, error)
	GetTransactionsBySubject(ctx context.Context, userId string, page, limit uint64, subject string, personOrBusiness string) ([]domain.Transaction, any, any, error)
	GetTransactionById(ctx context.Context, id string) (*domain.Transaction, error)
	CreateTransaction(ctx context.Context, createTransaction *domain.Transaction) (*domain.Transaction, error)
	UpdateTransaction(ctx context.Context, id string, updatedTransaction *domain.Transaction) (*domain.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
	DeleteTransactionsByUserId(ctx context.Context, id string) error
}

type TransactionService interface {
	GetStatus(ctx context.Context) string
	GetTransactionsByUserId(ctx context.Context, page, limit uint64, userId string) ([]domain.Transaction, any, any, error)
	GetTransactionsByDate(ctx context.Context, userId string, page, limit uint64, year int, month int) ([]domain.Transaction, any, any, error)
	GetTransactionsBySubject(ctx context.Context, userId string, page, limit uint64, subject string, personOrBusiness string) ([]domain.Transaction, any, any, error)
	GetTransactionById(ctx context.Context, id string) (*domain.Transaction, error)
	CreateTransaction(ctx context.Context, createTransaction *domain.Transaction) (*domain.Transaction, error)
	UpdateTransaction(ctx context.Context, id string, updatedTransaction *domain.Transaction) (*domain.Transaction, error)
	UpdateTotalOrigin(ctx context.Context, originId string, transactionType string, amount float64) error
	DeleteTransaction(ctx context.Context, id string) error
}
