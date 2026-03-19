package transactions

import "time"

type Transaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id,omitempty"`
	OccurredAt  string    `json:"occurred_at"` // YYYY-MM-DD
	Type        string    `json:"type"`        // income | expense
	AmountCents int64     `json:"amount_cents"`
	CategoryID  string    `json:"category_id"`
	Merchant    string    `json:"merchant"`
	Note        *string   `json:"note,omitempty"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListFilter struct {
	From       string
	To         string
	CategoryID string
	Q          string
	Limit      int
	Offset     int
}
