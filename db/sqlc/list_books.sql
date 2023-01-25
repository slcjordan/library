-- ListBooks returns a list of books.
-- name: ListBooks :many

SELECT isbn, title
FROM book
WHERE title > @page_token
ORDER BY title
LIMIT @total_size;
