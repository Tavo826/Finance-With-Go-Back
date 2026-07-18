package service

import (
	"context"
	"personal-finance/core/domain"
	"personal-finance/core/port"
	"time"
)

type ReportService struct {
	authService        port.AuthService
	transactionService port.TransactionService
	originService      port.OriginService
	mailAdapter        port.MailReportAdapter
}

func NewReportService(
	authService port.AuthService,
	transactionService port.TransactionService,
	originService port.OriginService,
	mailAdapter port.MailReportAdapter) *ReportService {

	return &ReportService{
		authService,
		transactionService,
		originService,
		mailAdapter,
	}
}

func (rs *ReportService) GenerateMonthlyReport(ctx context.Context, userId string) error {

	var report domain.Report
	var transactionList []domain.Transaction
	var page uint64 = 1
	var limit uint64 = 200

	var now = time.Now()
	var lastMonth = now.AddDate(0, -1, 0).Month()

	user, err := rs.authService.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	report.UserId = user.ID
	report.Username = user.Username
	report.UserEmail = user.Email
	report.Month = lastMonth
	report.Year = now.Year()

	origins, err := rs.originService.GetOriginsByUserId(ctx, user.ID)
	if err != nil {
		return err
	}

	report.NetBalance = calculateUserTotalNetwork(origins)

	for {
		transactions, _, totalPages, err := rs.transactionService.GetTransactionsByDate(ctx, user.ID, page, limit, now.Year(), int(lastMonth))
		if err != nil {
			return err
		}

		transactionList = append(transactionList, transactions...)

		page += 1

		if int(page) > totalPages {
			break
		}
	}

	filteredTransactions := filterTransactionsByType(transactionList)

	report.TotalIncome, report.TotalExpenses = calculateIncomeAndExpenses(filteredTransactions)
	report.OriginSummary = calculateOriginSummary(filteredTransactions, origins)
	report.CategorySummary = calculateCategorySummary(filteredTransactions)

	return rs.mailAdapter.SendMail(report)
}

func filterTransactionsByType(transactionList []domain.Transaction) []domain.Transaction {

	var filteredTransactionList []domain.Transaction

	for _, transaction := range transactionList {
		if transaction.Subject == "Payment" || transaction.Subject == "Expense" {
			filteredTransactionList = append(filteredTransactionList, transaction)
		}
	}

	return filteredTransactionList
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

func calculateCategorySummary(transactions []domain.Transaction) []domain.CategorySummary {

	totalMap := make(map[string]float64)
	countMap := make(map[string]int)
	var categoryOrder []string
	var categorySummaryList []domain.CategorySummary

	for _, transaction := range transactions {

		if transaction.OutputCategory == "" {
			continue
		}

		if _, exists := totalMap[transaction.OutputCategory]; !exists {
			categoryOrder = append(categoryOrder, transaction.OutputCategory)
		}

		totalMap[transaction.OutputCategory] += transaction.Amount
		countMap[transaction.OutputCategory] += 1
	}

	for _, category := range categoryOrder {
		var categorySummary domain.CategorySummary
		categorySummary.OutputCategory = category
		categorySummary.TotalExpenses = totalMap[category]
		categorySummary.Count = countMap[category]

		categorySummaryList = append(categorySummaryList, categorySummary)
	}

	return categorySummaryList
}
