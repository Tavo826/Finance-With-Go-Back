package service

import (
	"context"
	"errors"
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
) ([]domain.Transaction, int64, int, error) {

	transactions, totalDocuments, totalPages, err := ts.transactionRepo.GetTransactionsByDate(ctx, userId, page, limit, year, month)
	if err != nil {
		return nil, 0, 0, domain.ErrInternal
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
		if errors.Is(err, domain.ErrDataNotFound) {
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

	actualTransaction, err := ts.GetTransactionById(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := ts.reconcileOriginBalance(ctx, actualTransaction, transaction); err != nil {
		return nil, err
	}

	_, err = ts.transactionRepo.UpdateTransaction(ctx, id, transaction)
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

// reconcileOriginBalance applies the balance delta on the origin(s) affected
// by editing a transaction: origin change, type change and/or amount change.
func (ts *TransactionService) reconcileOriginBalance(ctx context.Context, actualTransaction *domain.Transaction, updatedTransaction *domain.Transaction) error {

	originId := ""
	transactionType := ""
	updateOrigin := false
	amount := float64(0)

	if *actualTransaction.OriginId != *updatedTransaction.OriginId {

		if actualTransaction.OriginId == nil {
			updateOrigin = true
			originId = *updatedTransaction.OriginId

			transactionType = updatedTransaction.Type
			amount = updatedTransaction.Amount
		} else {

			updateOrigin = true

			originId = *actualTransaction.OriginId

			if actualTransaction.Type == "Income" {
				transactionType = "Output"
			} else {
				transactionType = "Income"
			}

			amount = actualTransaction.Amount

			if err := ts.UpdateTotalOrigin(ctx, originId, transactionType, amount); err != nil {
				return err
			}

			originId = *updatedTransaction.OriginId

			transactionType = updatedTransaction.Type
			amount = updatedTransaction.Amount
		}

	} else if *actualTransaction.OriginId == *updatedTransaction.OriginId {

		originId = *updatedTransaction.OriginId
		transactionType = updatedTransaction.Type

		if actualTransaction.Type != updatedTransaction.Type {

			updateOrigin = true
			amount = actualTransaction.Amount + updatedTransaction.Amount
		} else {

			if actualTransaction.Amount != updatedTransaction.Amount {

				updateOrigin = true

				if actualTransaction.Amount > updatedTransaction.Amount {
					amount = actualTransaction.Amount - updatedTransaction.Amount

					if actualTransaction.Type == "Income" {
						transactionType = "Output"
					} else {
						transactionType = "Income"
					}
				} else {
					amount = updatedTransaction.Amount - actualTransaction.Amount
				}
			}
		}
	}

	if updateOrigin {
		return ts.UpdateTotalOrigin(ctx, originId, transactionType, amount)
	}

	return nil
}

func (ts *TransactionService) UpdateTotalOrigin(ctx context.Context, originId string, transactionType string, amount float64) error {

	origin, err := ts.originRepo.GetOriginById(ctx, originId)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		if errors.Is(err, domain.ErrDataNotFound) {
			return domain.ErrNoDocuments
		}
		return domain.ErrInternal
	}

	if transactionType == "Income" {
		origin.Total += amount
	} else if transactionType == "Output" {
		origin.Total -= amount
	}

	_, err = ts.originRepo.UpdateOrigin(ctx, originId, origin)
	if err != nil {
		return domain.ErrInternal
	}

	return nil
}

func (ts *TransactionService) DeleteTransaction(ctx context.Context, id string) error {

	return ts.transactionRepo.DeleteTransaction(ctx, id)
}
