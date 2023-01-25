-- DeleteBook deletes a single book.
-- name: DeleteBook :exec

DELETE FROM book WHERE isbn = @isbn;
