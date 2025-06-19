package service

import (
	"context"
	"personal-finance/core/domain"
	"personal-finance/core/port"
)

type TransactionService struct {
	transactionRepo port.TransactionRepository
	originRepo      port.OriginRepository
}

func NewTransactionService(transactionRepo port.TransactionRepository, originRepo port.OriginRepository) *TransactionService {

	return &TransactionService{
		transactionRepo,
		originRepo,
	}
}

func (ts *TransactionService) GetStatus(ctx context.Context) string {

	return "OK"
}

func (ts *TransactionService) GetTransactionsByUserId(ctx context.Context, page, limit uint64, userId string) ([]domain.Transaction, any, any, error) {

	transactions, totalDocuments, totalPages, err := ts.transactionRepo.GetTransactionsByUserId(ctx, page, limit, userId)
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

	transactions, totalDocuments, totalPages, err := ts.transactionRepo.GetTransactionsByDate(ctx, userId, page, limit, year, month)
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

	transactions, totalDocuments, totalPages, err := ts.transactionRepo.GetTransactionsBySubject(ctx, userId, page, limit, subject, personOrBusiness)
	if err != nil {
		return nil, nil, nil, domain.ErrInternal
	}

	return transactions, totalDocuments, totalPages, nil
}

func (ts *TransactionService) GetTransactionById(ctx context.Context, id string) (*domain.Transaction, error) {

	transaction, err := ts.transactionRepo.GetTransactionById(ctx, id)
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

	transaction, err := ts.transactionRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return transaction, nil
}

func (ts *TransactionService) UpdateTransaction(ctx context.Context, id string, transaction *domain.Transaction) (*domain.Transaction, error) {

	_, err := ts.transactionRepo.UpdateTransaction(ctx, id, transaction)
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

func (ts *TransactionService) UpdateTotalOrigin(ctx context.Context, transaction *domain.Transaction, transactionType string) error {

	origin, err := ts.originRepo.GetOriginById(ctx, *transaction.OriginId)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		if err.Error() == domain.ErrNoDocuments.Error() {
			return domain.ErrNoDocuments
		}
		return domain.ErrInternal
	}

	if transactionType == "Input" {
		origin.Total += transaction.Amount
	} else if transactionType == "Output" {
		origin.Total -= transaction.Amount
	}

	_, err = ts.originRepo.UpdateOrigin(ctx, *transaction.OriginId, origin)
	if err != nil {
		return domain.ErrInternal
	}

	return nil
}

func (ts *TransactionService) DeleteTransaction(ctx context.Context, id string) error {

	return ts.transactionRepo.DeleteTransaction(ctx, id)
}
