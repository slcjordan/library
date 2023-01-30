-- CreateBook creates a single book.
-- name: CreateBook :exec

INSERT INTO book (isbn, title)
VALUES (@isbn, @title);
