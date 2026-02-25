package categories

import "time"

type Category struct {
	ID        string
	UserID    *string
	Name      string
	Type      string
	IsSystem  bool
	CreatedAt time.Time
}
