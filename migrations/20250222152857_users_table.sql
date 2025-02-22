-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   telegram_id VARCHAR(255) UNIQUE NOT NULL,
   created_at TIMESTAMP DEFAULT now() NOT NULL,
   updated_at TIMESTAMP DEFAULT now() NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
