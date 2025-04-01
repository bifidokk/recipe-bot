package entity

import (
	"database/sql"
	"strings"
	"time"
)

type Recipe struct {
	ID                 int            `db:"id"`
	Title              string         `db:"title"`
	Body               string         `db:"body"`
	RecipeMarkdownText string         `db:"markdown"`
	Source             string         `db:"source"`
	SourceID           string         `db:"source_id"`
	SourceIDType       string         `db:"source_id_type"`
	AudioURL           string         `db:"audio_url"`
	ShareURL           sql.NullString `db:"share_url"`
	CoverFileID        sql.NullString `db:"cover_file_id"`
	UserID             int            `db:"user_id"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
}

func (r *Recipe) GetRecipeMarkdownView() string {
	return escapeMarkdownV2(r.RecipeMarkdownText)
}

func escapeMarkdownV2(text string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}
