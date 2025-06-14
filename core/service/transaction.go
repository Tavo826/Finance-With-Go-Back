package service

import (
	"context"
	"personal-finance/core/domain"
	"personal-finance/core/port"
)

type TransactionService struct {
	repo port.TransactionRepository
}

func NewTransactionService(repo port.TransactionRepository) *TransactionService {

	return &TransactionService{
		repo,
	}
}

func (ts *TransactionService) GetStatus(ctx context.Context) string {

	return "OK"
}

func (ts *TransactionService) GetTransactionsByUserId(ctx context.Context, page, limit uint64, userId string) ([]domain.Transaction, any, any, error) {

	transactions, totalDocuments, totalPages, err := ts.repo.GetTransactionsByUserId(ctx, page, limit, userId)
	if err != nil {
		return nil, nil, nil, domain.ErrInternal
	}

	return transactions, totalDocuments, totalPages, nil
}

func (ts *TransactionService) GetTransactionsByDate(
	ctx context.Context,
	userId string,
	page, limit uint64,
	year int,
	month int,
) ([]domain.Transaction, any, any, error) {

	transactions, totalDocuments, totalPages, err := ts.repo.GetTransactionsByDate(ctx, userId, page, limit, year, month)
	if err != nil {
		return nil, nil, nil, domain.ErrInternal
	}

	return transactions, totalDocuments, totalPages, nil
}

func (ts *TransactionService) GetTransactionsBySubject(
	ctx context.Context,
	userId string,
	page, limit uint64,
	subject string,
	personOrBusiness string,
) ([]domain.Transaction, any, any, error) {

	transactions, totalDocuments, totalPages, err := ts.repo.GetTransactionsBySubject(ctx, userId, page, limit, subject, personOrBusiness)
	if err != nil {
		return nil, nil, nil, domain.ErrInternal
	}

	return transactions, totalDocuments, totalPages, nil
}

func (ts *TransactionService) GetTransaction(ctx context.Context, id string) (*domain.Transaction, error) {

	transaction, err := ts.repo.GetTransaction(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		if err.Error() == domain.ErrNoDocuments.Error() {
			return nil, domain.ErrNoDocuments
		}
		return nil, domain.ErrInternal
	}

	return transaction, nil
}

func (ts *TransactionService) CreateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {

	transaction, err := ts.repo.CreateTransaction(ctx, transaction)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return transaction, nil
}

func (ts *TransactionService) UpdateTransaction(ctx context.Context, id string, transaction *domain.Transaction) (*domain.Transaction, error) {

	_, err := ts.repo.UpdateTransaction(ctx, id, transaction)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		if err == domain.ErrDataNotFound {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.ErrInternal
	}

	return transaction, nil
}

func (ts *TransactionService) DeleteTransaction(ctx context.Context, id string) error {

	return ts.repo.DeleteTransaction(ctx, id)
}
