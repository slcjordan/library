-- GetBook fetches a single book.
-- name: GetBook :one

SELECT title, isbn FROM book WHERE isbn = @isbn;
