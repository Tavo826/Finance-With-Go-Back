package port

import (
	"context"
	"personal-finance/core/domain"
)

type ReportService interface {
	GenerateMonthlyReport(ctx context.Context, userId string) error
}

type MailReportAdapter interface {
	SendMail(report domain.Report) error
}
