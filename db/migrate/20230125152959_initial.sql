-- +goose Up
-- +goose StatementBegin
CREATE TABLE book (
  isbn BIGINT NOT NULL PRIMARY KEY,
  title TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS book;
-- +goose StatementEnd
