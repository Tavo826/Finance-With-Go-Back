package domain

import "time"

type Report struct {
	UserId        string          `json:"user_id"`
	Username      string          `json:"username"`
	UserEmail     string          `json:"email"`
	Month         time.Month      `json:"month"`
	Year          int             `json:"year"`
	TotalIncome   float64         `json:"total_income"`
	TotalExpenses float64         `json:"total_expenses"`
	NetBalance    float64         `json:"net_balance"`
	OriginSummary []OriginSummary `json:"origin_summary"`
}

type OriginSummary struct {
	OriginName    string  `json:"origin_name"`
	TotalIncome   float64 `json:"total_income"`
	TotalExpenses float64 `json:"total_expenses"`
	OriginBalance float64 `json:"origin_balance"`
}
