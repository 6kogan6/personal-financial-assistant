package reports

type Summary struct {
	Month  string `json:"month"`
	Totals Totals `json:"totals"`

	ByCategory []ByCategoryItem `json:"by_category"`
	Daily      []DailyItem      `json:"daily"`
}

type Totals struct {
	IncomeCents  int64 `json:"income_cents"`
	ExpenseCents int64 `json:"expense_cents"`
	BalanceCents int64 `json:"balance_cents"`
}

type ByCategoryItem struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	Type         string `json:"type"` // expense|income
	AmountCents  int64  `json:"amount_cents"`
}

type DailyItem struct {
	Date         string `json:"date"` // YYYY-MM-DD
	IncomeCents  int64  `json:"income_cents"`
	ExpenseCents int64  `json:"expense_cents"`
}
