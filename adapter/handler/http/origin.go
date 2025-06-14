package http

import (
	"log"
	"personal-finance/adapter/handler/http/dto"
	"personal-finance/core/domain"
	"personal-finance/core/port"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OriginHandler struct {
	service  port.OriginService
	validate *validator.Validate
}

func NewOriginHandler(service port.OriginService, validate *validator.Validate) *OriginHandler {
	return &OriginHandler{
		service,
		validate,
	}
}

func (oh *OriginHandler) GetOriginsByUserId(ctx *gin.Context) {

	var req dto.OriginByUserId
	var originList []dto.OriginResponse

	log.Println("GetOriginsByUserId")

	if err := ctx.Bind(&req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	origins, err := oh.service.GetOriginsByUserId(ctx, req.UserId)
	if err != nil {
		log.Println("Handle error: ", err)
		dto.HandleError(ctx, err)
		return
	}

	for _, origin := range origins {
		originList = append(originList, dto.NewOriginResponse(&origin))
	}

	if originList == nil {
		originList = []dto.OriginResponse{}
	}

	dto.HandleSuccess(ctx, originList)
}

func (oh *OriginHandler) GetOriginById(ctx *gin.Context) {

	var request dto.IdRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	origin, err := oh.service.GetOriginById(ctx, request.ID)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	response := dto.NewOriginResponse(origin)

	dto.HandleSuccess(ctx, response)
}

func (oh *OriginHandler) CreateOrigin(ctx *gin.Context) {

	var req dto.OriginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	if err := oh.validate.Struct(req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	origin := domain.Origin{
		UserId:    req.UserId,
		Name:      req.Name,
		Total:     req.Total,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := oh.service.CreateOrigin(ctx, &origin)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	response := dto.NewOriginResponse(&origin)

	dto.HandleSuccess(ctx, response)
}

func (oh *OriginHandler) UpdateOrigin(ctx *gin.Context) {

	var req dto.OriginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	if err := oh.validate.Struct(req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	id := ctx.Param("id")

	origin := domain.Origin{
		UserId:    req.UserId,
		Name:      req.Name,
		Total:     req.Total,
		CreatedAt: req.CreatedAt,
		UpdatedAt: time.Now(),
	}

	_, err := oh.service.UpdateOrigin(ctx, id, &origin)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	response := dto.NewOriginResponse(&origin)

	dto.HandleSuccess(ctx, response)
}

func (oh *OriginHandler) DeleteOrigin(ctx *gin.Context) {

	var request dto.IdRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	err := oh.service.DeleteOrigin(ctx, request.ID)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	dto.HandleSuccess(ctx, nil)
}
