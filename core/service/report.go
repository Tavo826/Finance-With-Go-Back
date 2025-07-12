package service

import (
	"personal-finance/core/domain"
	"personal-finance/core/port"
)

type ReportService struct {
	mailAdapter port.MailReportAdapter
}

func NewReportService(mailAdapter port.MailReportAdapter) *ReportService {

	return &ReportService{
		mailAdapter,
	}
}

func (rs *ReportService) SendReport(report domain.Report) error {

	if err := rs.mailAdapter.SendMail(report); err != nil {
		return err
	}

	return nil
}
