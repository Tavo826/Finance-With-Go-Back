package http

import (
	"personal-finance/adapter/handler/http/dto"
	"personal-finance/core/domain"
	"personal-finance/core/port"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TransactionHandler struct {
	service  port.TransactionService
	validate *validator.Validate
}

func NewTransactionHandler(service port.TransactionService, validate *validator.Validate) *TransactionHandler {
	return &TransactionHandler{
		service,
		validate,
	}
}

func (th *TransactionHandler) GetStatus(ctx *gin.Context) {

	response := th.service.GetStatus(ctx)
	dto.HandleSuccess(ctx, response)
}

func (th *TransactionHandler) GetTransactionsByUserId(ctx *gin.Context) {

	var req dto.TransactionByUserRequest
	var transactionList []dto.TransactionResponse

	if err := ctx.Bind(&req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	transactions, totalDocuments, totalPages, err := th.service.GetTransactionsByUserId(ctx, req.Page, req.Limit, req.UserId)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	for _, transaction := range transactions {
		transactionList = append(transactionList, dto.NewTransactionResponse(&transaction))
	}

	if transactionList == nil {
		transactionList = []dto.TransactionResponse{}
	}

	response := dto.NewPaginatedResponse(
		req.Page,
		req.Limit,
		totalDocuments.(int64),
		totalPages.(int),
		transactionList,
	)

	dto.HandleSuccess(ctx, response)
}

func (th *TransactionHandler) GetTransactionsByDate(ctx *gin.Context) {

	var req dto.DateFilterRequest
	var transactionList []dto.TransactionResponse

	if err := ctx.Bind(&req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	transactions, totalDocuments, totalPages, err := th.service.GetTransactionsByDate(ctx, req.UserId, req.Page, req.Limit, req.Year, req.Month)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	for _, transaction := range transactions {
		transactionList = append(transactionList, dto.NewTransactionResponse(&transaction))
	}

	if transactionList == nil {
		transactionList = []dto.TransactionResponse{}
	}

	response := dto.NewPaginatedResponse(
		req.Page,
		req.Limit,
		totalDocuments.(int64),
		totalPages.(int),
		transactionList,
	)

	dto.HandleSuccess(ctx, response)
}

func (th *TransactionHandler) GetTransactionsBySubject(ctx *gin.Context) {

	var req dto.SubjectFilterRequest
	var transactionList []dto.TransactionResponse

	if err := ctx.Bind(&req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	transactions, totalDocuments, totalPages, err := th.service.GetTransactionsBySubject(ctx, req.UserId, req.Page, req.Limit, req.Subject, req.PersonOrBusiness)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	for _, transaction := range transactions {
		transactionList = append(transactionList, dto.NewTransactionResponse(&transaction))
	}

	if transactionList == nil {
		transactionList = []dto.TransactionResponse{}
	}

	response := dto.NewPaginatedResponse(
		req.Page,
		req.Limit,
		totalDocuments.(int64),
		totalPages.(int),
		transactionList,
	)

	dto.HandleSuccess(ctx, response)
}

func (th *TransactionHandler) GetTransactionById(ctx *gin.Context) {
	var request dto.IdRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	transaction, err := th.service.GetTransactionById(ctx, request.ID)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	response := dto.NewTransactionResponse(transaction)

	dto.HandleSuccess(ctx, response)
}

func (th *TransactionHandler) CreateTransaction(ctx *gin.Context) {

	var req dto.TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	if err := th.validate.Struct(req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	transaction := domain.Transaction{
		Amount:           req.Amount,
		UserId:           req.UserId,
		OriginId:         &req.OriginId,
		Type:             req.Type,
		Subject:          req.Subject,
		PersonOrBusiness: req.PersonOrBusiness,
		Description:      req.Description,
		CreatedAtString:  req.CreatedAtString,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, err := th.service.CreateTransaction(ctx, &transaction)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	if *transaction.OriginId != "" {
		err = th.service.UpdateTotalOrigin(ctx, *transaction.OriginId, transaction.Type, transaction.Amount)
		if err != nil {
			dto.HandleError(ctx, err)
			return
		}
	}

	response := dto.NewTransactionResponse(&transaction)

	dto.HandleSuccess(ctx, response)
}

func (th *TransactionHandler) UpdateTransaction(ctx *gin.Context) {

	var req dto.TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	if err := th.validate.Struct(req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	id := ctx.Param("id")

	updatedTransaction := domain.Transaction{
		Amount:           req.Amount,
		UserId:           req.UserId,
		OriginId:         &req.OriginId,
		Type:             req.Type,
		Subject:          req.Subject,
		PersonOrBusiness: req.PersonOrBusiness,
		Description:      req.Description,
		CreatedAtString:  req.CreatedAtString,
		CreatedAt:        req.CreatedAt,
		UpdatedAt:        time.Now(),
	}

	actualTransaction, err := th.service.GetTransactionById(ctx, id)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	verifyAndUpdateOrigin(ctx, th, actualTransaction, updatedTransaction)

	_, err = th.service.UpdateTransaction(ctx, id, &updatedTransaction)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	response := dto.NewTransactionResponse(&updatedTransaction)

	dto.HandleSuccess(ctx, response)
}

func verifyAndUpdateOrigin(ctx *gin.Context, th *TransactionHandler, actualTransaction *domain.Transaction, updatedTransaction domain.Transaction) {

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

			err := th.service.UpdateTotalOrigin(ctx, originId, transactionType, float64(amount))
			if err != nil {
				dto.HandleError(ctx, err)
				return
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

		err := th.service.UpdateTotalOrigin(ctx, originId, transactionType, float64(amount))
		if err != nil {
			dto.HandleError(ctx, err)
			return
		}
	}
}

func (th *TransactionHandler) DeleteTransaction(ctx *gin.Context) {

	var transactionType string = "Income"

	var request dto.IdRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	transaction, err := th.service.GetTransactionById(ctx, request.ID)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	if *transaction.OriginId != "" {

		if transaction.Type == "Income" {
			transactionType = "Output"
		}

		err = th.service.UpdateTotalOrigin(ctx, *transaction.OriginId, transactionType, transaction.Amount)
		if err != nil {
			dto.HandleError(ctx, err)
			return
		}
	}

	err = th.service.DeleteTransaction(ctx, request.ID)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	dto.HandleSuccess(ctx, nil)
}
