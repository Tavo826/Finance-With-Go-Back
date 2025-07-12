package http

import (
	"context"
	"personal-finance/core/domain"
	"personal-finance/core/port"
	"time"
)

type ReportHandler struct {
	userService        port.AuthService
	transactionService port.TransactionService
	originService      port.OriginService
	reportService      port.ReportService
}

func NewReportHandler(
	userService port.AuthService,
	transactionService port.TransactionService,
	originService port.OriginService,
	reportService port.ReportService) *ReportHandler {
	return &ReportHandler{
		userService,
		transactionService,
		originService,
		reportService,
	}
}

func (rh *ReportHandler) GenerateMonthlyTransactionReport(ctx context.Context) error {

	var report domain.Report
	var transactionList []domain.Transaction
	var origins []domain.Origin
	var page uint64 = 1
	var limit uint64 = 200

	var now = time.Now()
	var lastMonth = now.AddDate(0, -1, 0).Month()

	users, err := rh.userService.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {

		report.UserId = user.ID
		report.Username = user.Username
		report.UserEmail = user.Email
		report.Month = lastMonth
		report.Year = now.Year()

		origins, err = rh.originService.GetOriginsByUserId(ctx, user.ID)
		if err != nil {
			return err
		}

		report.NetBalance = calculateUserTotalNetwork(origins)

		for {
			transactions, _, totalPages, err := rh.transactionService.GetTransactionsByDate(ctx, user.ID, page, limit, now.Year(), int(lastMonth))
			if err != nil {
				return err
			}

			transactionList = append(transactionList, transactions...)

			page += 1

			if int(page) > totalPages.(int) {
				break
			}
		}

		report.TotalIncome, report.TotalExpenses = calculateIncomeAndExpenses(transactionList)
		report.OriginSummary = calculateOriginSummary(transactionList, origins)

		if err = rh.reportService.SendReport(report); err != nil {
			return err
		}
	}

	return nil
}

func calculateUserTotalNetwork(origins []domain.Origin) float64 {

	var totalNetwork float64 = 0

	for _, origin := range origins {
		totalNetwork += origin.Total
	}

	return totalNetwork
}

func calculateIncomeAndExpenses(transactions []domain.Transaction) (float64, float64) {

	var totalIncome float64 = 0
	var totalExpenses float64 = 0

	for _, transaction := range transactions {

		if transaction.Type == "Income" {
			totalIncome += transaction.Amount
		} else {
			totalExpenses += transaction.Amount
		}
	}

	return totalIncome, totalExpenses
}

func calculateOriginSummary(transactions []domain.Transaction, origins []domain.Origin) []domain.OriginSummary {

	incomeMap := make(map[string]float64)
	outputMap := make(map[string]float64)
	var originSummaryList []domain.OriginSummary

	for _, origin := range origins {
		incomeMap[origin.ID] = 0
		outputMap[origin.ID] = 0
	}

	for _, transaction := range transactions {

		if transaction.OriginId == nil {
			continue
		}

		originId := *transaction.OriginId

		if transaction.Type == "Income" {
			incomeMap[originId] += transaction.Amount
		} else {
			outputMap[originId] += transaction.Amount
		}
	}

	for _, origin := range origins {
		var originSummary domain.OriginSummary
		originSummary.OriginName = origin.Name
		originSummary.OriginBalance = origin.Total
		originSummary.TotalIncome = incomeMap[origin.ID]
		originSummary.TotalExpenses = outputMap[origin.ID]

		originSummaryList = append(originSummaryList, originSummary)
	}

	return originSummaryList
}
