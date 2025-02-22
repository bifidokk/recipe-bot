package entity

import "time"

type User struct {
	ID         int       `db:"id"`
	Name       string    `db:"name"`
	TelegramID string    `db:"telegram_id"`
	Language   string    `db:"language_code"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
