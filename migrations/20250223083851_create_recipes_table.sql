-- +goose Up
-- +goose StatementBegin
CREATE TABLE recipes (
     id SERIAL PRIMARY KEY,
     title VARCHAR(255) NOT NULL,
     body TEXT NOT NULL,
     markdown TEXT NOT NULL,
     source VARCHAR(255),
     source_link VARCHAR(512),
     audio_link VARCHAR(512),
     user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
     created_at TIMESTAMP DEFAULT now() NOT NULL,
     updated_at TIMESTAMP DEFAULT now() NOT NULL
);

CREATE INDEX idx_recipes_user_id ON recipes(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS recipes;
-- +goose StatementEnd
