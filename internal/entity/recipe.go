package entity

import "time"

type Recipe struct {
	ID                 int       `db:"id"`
	Title              string    `db:"title"`
	Body               string    `db:"body"`
	RecipeMarkdownText string    `db:"markdown"`
	Source             string    `db:"source"`
	SourceID           string    `db:"source_id"`
	SourceIDType       string    `db:"source_id_type"`
	AudioURL           string    `db:"audio_url"`
	UserID             int       `db:"user_id"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}
