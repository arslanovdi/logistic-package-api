-- +goose Up
CREATE TABLE IF NOT EXISTS package (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE package;
