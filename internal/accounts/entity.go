package accounts

import "time"

type Account struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
