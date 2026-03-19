package transactions

import "time"

type Transaction struct {
	ID          string
	UserID      string
	OccurredAt  string // YYYY-MM-DD
	Type        string
	AmountCents int64
	CategoryID  string
	Merchant    string
	Note        *string
	Source      string
	CreatedAt   time.Time
}

type ListFilter struct {
	From       string // YYYY-MM-DD optional
	To         string // YYYY-MM-DD optional
	CategoryID string // optional
	Q          string // optional search in merchant/note
}
