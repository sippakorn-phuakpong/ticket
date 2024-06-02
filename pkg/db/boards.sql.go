// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: boards.sql

package db

import (
	"context"

	"github.com/guregu/null"
)

const countBoardByUserID = `-- name: CountBoardByUserID :one
SELECT
  COUNT(*)
FROM
  boards
WHERE
  user_id = ?
`

func (q *Queries) CountBoardByUserID(ctx context.Context, userID uint64) (int64, error) {
	row := q.db.QueryRowContext(ctx, countBoardByUserID, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createBoard = `-- name: CreateBoard :exec
INSERT INTO
  boards (user_id, title, sort_order, created_at)
VALUES
  (?, ?, ?, NOW())
`

type CreateBoardParams struct {
	UserID    uint64      `db:"user_id" json:"user_id"`
	Title     null.String `db:"title" json:"title"`
	SortOrder uint32      `db:"sort_order" json:"sort_order"`
}

func (q *Queries) CreateBoard(ctx context.Context, arg CreateBoardParams) error {
	_, err := q.db.ExecContext(ctx, createBoard, arg.UserID, arg.Title, arg.SortOrder)
	return err
}

const deleteBoard = `-- name: DeleteBoard :exec
DELETE FROM
  boards
WHERE
  id = ?
`

func (q *Queries) DeleteBoard(ctx context.Context, id uint32) error {
	_, err := q.db.ExecContext(ctx, deleteBoard, id)
	return err
}

const getBoardByID = `-- name: GetBoardByID :one
SELECT
  id, user_id, title, sort_order, created_at, updated_at
FROM
  boards
WHERE
  id = ?
`

func (q *Queries) GetBoardByID(ctx context.Context, id uint32) (Board, error) {
	row := q.db.QueryRowContext(ctx, getBoardByID, id)
	var i Board
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.SortOrder,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getBoardsByUserID = `-- name: GetBoardsByUserID :many
SELECT
  id, user_id, title, sort_order, created_at, updated_at
FROM
  boards
WHERE
  user_id = ?
`

func (q *Queries) GetBoardsByUserID(ctx context.Context, userID uint64) ([]Board, error) {
	rows, err := q.db.QueryContext(ctx, getBoardsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Board{}
	for rows.Next() {
		var i Board
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.SortOrder,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLastInsertBoardViewByUserID = `-- name: GetLastInsertBoardViewByUserID :one
SELECT
  id, user_id, title, sort_order, created_at, updated_at, statuses
FROM
  board_view
WHERE
  board_view.user_id = ?
  AND id = (
    SELECT
      LAST_INSERT_ID()
    FROM
      boards
    LIMIT
      1
  )
`

func (q *Queries) GetLastInsertBoardViewByUserID(ctx context.Context, userID uint64) (BoardView, error) {
	row := q.db.QueryRowContext(ctx, getLastInsertBoardViewByUserID, userID)
	var i BoardView
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.SortOrder,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Statuses,
	)
	return i, err
}

const listBoardViewByUserID = `-- name: ListBoardViewByUserID :many
SELECT
  id, user_id, title, sort_order, created_at, updated_at, statuses
FROM
  board_view
WHERE
  user_id = ?
ORDER BY
  sort_order ASC
`

func (q *Queries) ListBoardViewByUserID(ctx context.Context, userID uint64) ([]BoardView, error) {
	rows, err := q.db.QueryContext(ctx, listBoardViewByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []BoardView{}
	for rows.Next() {
		var i BoardView
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.SortOrder,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Statuses,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBoard = `-- name: UpdateBoard :exec
UPDATE
  boards
SET
  title = ?,
  updated_at = NOW()
WHERE
  id = ?
`

type UpdateBoardParams struct {
	Title null.String `db:"title" json:"title"`
	ID    uint32      `db:"id" json:"id"`
}

func (q *Queries) UpdateBoard(ctx context.Context, arg UpdateBoardParams) error {
	_, err := q.db.ExecContext(ctx, updateBoard, arg.Title, arg.ID)
	return err
}
