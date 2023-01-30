-- UpdateBook updates a single book.
-- name: UpdateBook :exec

UPDATE book SET title = @title WHERE isbn = @isbn;
