package http

import (
	"personal-finance/adapter/handler/http/dto"
	"personal-finance/core/port"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService port.ReportService
}

func NewReportHandler(reportService port.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService,
	}
}

func (rh *ReportHandler) GenerateMonthlyTransactionReport(ctx *gin.Context) {

	var request dto.RequestByUserId
	if err := ctx.Bind(&request); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	if err := rh.reportService.GenerateMonthlyReport(ctx, request.UserId); err != nil {
		dto.HandleError(ctx, err)
		return
	}

	dto.HandleSuccess(ctx, "Mail sended successfully")
}
