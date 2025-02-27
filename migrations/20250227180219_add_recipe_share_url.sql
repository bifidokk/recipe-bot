-- +goose Up
-- +goose StatementBegin
ALTER TABLE recipes ADD COLUMN share_url TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE recipes DROP COLUMN share_url;
-- +goose StatementEnd
