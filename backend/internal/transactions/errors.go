package transactions

import "errors"

var (
	ErrInvalidUserID        = errors.New("invalid user id")
	ErrInvalidTransactionID = errors.New("invalid transaction id")
	ErrInvalidCategoryID    = errors.New("invalid category id")
	ErrInvalidType          = errors.New("invalid transaction type")
	ErrInvalidAmount        = errors.New("invalid amount")
	ErrInvalidMerchant      = errors.New("invalid merchant")
	ErrInvalidOccurredAt    = errors.New("invalid occurred_at")
	ErrInvalidFromDate      = errors.New("invalid from")
	ErrInvalidToDate        = errors.New("invalid to")
	ErrTransactionNotFound  = errors.New("transaction not found")
	ErrCategoryNotAllowed   = errors.New("category not allowed")
)
