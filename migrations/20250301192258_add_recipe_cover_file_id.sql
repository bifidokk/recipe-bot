-- +goose Up
-- +goose StatementBegin
ALTER TABLE recipes ADD COLUMN cover_file_id VARCHAR(255) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE recipes DROP COLUMN cover_file_id;
-- +goose StatementEnd
