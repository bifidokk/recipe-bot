-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN recipe_limit INT DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN recipe_limit;
-- +goose StatementEnd
