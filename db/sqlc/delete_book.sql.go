// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: delete_book.sql

package sqlc

import (
	"context"
)

const deleteBook = `-- name: DeleteBook :exec

DELETE FROM book WHERE isbn = $1
`

// DeleteBook deletes a single book.
func (q *Queries) DeleteBook(ctx context.Context, isbn int64) error {
	_, err := q.db.Exec(ctx, deleteBook, isbn)
	return err
}