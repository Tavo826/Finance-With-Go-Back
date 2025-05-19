package dto

import (
	"net/http"
	"time"

	"personal-finance/core/domain"

	"github.com/gin-gonic/gin"
)

type TransactionResponse struct {
	ID               any       `json:"_id"`
	UserId           string    `json:"user_id"`
	Amount           float64   `json:"amount"`
	Type             string    `json:"type"`
	Subject          string    `json:"subject"`
	PersonOrBusiness string    `json:"person_business"`
	Description      string    `json:"description"`
	CreatedAtString  string    `json:"created"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Body    any    `json:"body"`
}

type PaginatedResponse struct {
	Page           uint64                `json:"page"`
	Limit          uint64                `json:"limit"`
	TotalDocuments int64                 `json:"totalDocuments"`
	TotalPages     int                   `json:"totalPages"`
	Data           []TransactionResponse `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var errorStatusMap = map[error]int{
	domain.ErrInternal:                   http.StatusInternalServerError,
	domain.ErrDataNotFound:               http.StatusNotFound,
	domain.ErrNoDocuments:                http.StatusNotFound,
	domain.ErrConflictingData:            http.StatusConflict,
	domain.ErrInvalidCredentials:         http.StatusUnauthorized,
	domain.ErrUnauthorized:               http.StatusUnauthorized,
	domain.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	domain.ErrInvalidToken:               http.StatusUnauthorized,
	domain.ErrExpiredToken:               http.StatusUnauthorized,
	domain.ErrForbidden:                  http.StatusForbidden,
	domain.ErrUserAlreadyExists:          http.StatusBadRequest,
	domain.ErrNoUpdatedData:              http.StatusBadRequest,
}

func NewTransactionResponse(transaction *domain.Transaction) TransactionResponse {

	return TransactionResponse{
		ID:               transaction.ID,
		UserId:           transaction.UserId,
		Amount:           transaction.Amount,
		Type:             transaction.Type,
		Subject:          transaction.Subject,
		PersonOrBusiness: transaction.PersonOrBusiness,
		Description:      transaction.Description,
		CreatedAtString:  transaction.CreatedAtString,
		CreatedAt:        transaction.CreatedAt,
		UpdatedAt:        transaction.UpdatedAt,
	}
}

func NewPaginatedResponse(
	page uint64,
	limit uint64,
	totalDocuments int64,
	totalPages int,
	transactionList []TransactionResponse,
) PaginatedResponse {

	return PaginatedResponse{
		Page:           page,
		Limit:          limit,
		TotalDocuments: totalDocuments,
		TotalPages:     totalPages,
		Data:           transactionList,
	}
}

func ValidationError(ctx *gin.Context, err error) {

	errorResponse := newErrorResponse(err.Error())
	ctx.JSON(http.StatusBadRequest, errorResponse)
}

func HandleSuccess(ctx *gin.Context, data any) {
	response := newResponse(true, "Success", data)
	ctx.JSON(http.StatusOK, response)
}

func newResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
		Body:    data,
	}
}

func HandleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]

	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errorResponse := newErrorResponse(err.Error())
	ctx.JSON(statusCode, errorResponse)
}

func newErrorResponse(errorMessage string) ErrorResponse {

	return ErrorResponse{
		Success: false,
		Message: errorMessage,
	}
}
