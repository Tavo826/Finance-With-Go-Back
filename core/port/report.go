package port

import "personal-finance/core/domain"

type ReportService interface {
	SendReport(report domain.Report) error
}

type MailReportAdapter interface {
	SendMail(report domain.Report) error
}
