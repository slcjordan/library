// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: list_books.sql

package sqlc

import (
	"context"
)

const listBooks = `-- name: ListBooks :many

SELECT isbn, title
FROM book
WHERE title > $1
ORDER BY title
LIMIT $2
`

type ListBooksParams struct {
	PageToken string
	TotalSize int32
}

// ListBooks returns a list of books.
func (q *Queries) ListBooks(ctx context.Context, arg ListBooksParams) ([]Book, error) {
	rows, err := q.db.Query(ctx, listBooks, arg.PageToken, arg.TotalSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Book
	for rows.Next() {
		var i Book
		if err := rows.Scan(&i.Isbn, &i.Title); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
