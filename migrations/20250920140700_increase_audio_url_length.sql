-- +goose Up
-- +goose StatementBegin
ALTER TABLE recipes ALTER COLUMN audio_url TYPE TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE recipes ALTER COLUMN audio_url TYPE VARCHAR(512);
-- +goose StatementEnd
