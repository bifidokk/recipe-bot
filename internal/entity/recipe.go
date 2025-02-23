package entity

import "time"

type Recipe struct {
	ID                 int       `db:"id"`
	Title              string    `db:"title"`
	Body               string    `db:"body"`
	RecipeMarkdownText string    `db:"markdown"`
	Source             string    `db:"source"`
	SourceLink         string    `db:"source_link"`
	AudioLink          string    `db:"audio_link"`
	UserID             int       `db:"user_id"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}
