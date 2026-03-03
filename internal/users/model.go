package users

import "time"

type User struct {
	ID        string
	AccountID string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
