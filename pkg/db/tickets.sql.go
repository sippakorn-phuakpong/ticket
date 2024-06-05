// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: tickets.sql

package db

import (
	"context"
	"strings"

	null "github.com/guregu/null/v5"
)

const countTicketByStatusID = `-- name: CountTicketByStatusID :one
SELECT
  COUNT(*)
FROM
  tickets
WHERE
  status_id = ?
`

func (q *Queries) CountTicketByStatusID(ctx context.Context, statusID uint32) (int64, error) {
	row := q.db.QueryRowContext(ctx, countTicketByStatusID, statusID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countTicketWithBoard = `-- name: CountTicketWithBoard :one
SELECT
  COUNT(tickets.id)
FROM
  tickets
  JOIN statuses ON tickets.status_id = statuses.id
  JOIN boards ON statuses.board_id = boards.id
WHERE
  tickets.id IN (/*SLICE:ids*/?)
  AND statuses.board_id = ?
  AND boards.user_id = ?
`

type CountTicketWithBoardParams struct {
	Ids     []uint64 `db:"ids" json:"ids"`
	BoardID uint32   `db:"board_id" json:"board_id"`
	UserID  uint64   `db:"user_id" json:"user_id"`
}

func (q *Queries) CountTicketWithBoard(ctx context.Context, arg CountTicketWithBoardParams) (int64, error) {
	query := countTicketWithBoard
	var queryParams []interface{}
	if len(arg.Ids) > 0 {
		for _, v := range arg.Ids {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:ids*/?", strings.Repeat(",?", len(arg.Ids))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:ids*/?", "NULL", 1)
	}
	queryParams = append(queryParams, arg.BoardID)
	queryParams = append(queryParams, arg.UserID)
	row := q.db.QueryRowContext(ctx, query, queryParams...)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createTicket = `-- name: CreateTicket :exec
INSERT INTO
  tickets (
    status_id,
    title,
    description,
    contact,
    sort_order,
    created_at
  )
VALUES
  (?, ?, ?, ?, ?, NOW())
`

type CreateTicketParams struct {
	StatusID    uint32      `db:"status_id" json:"status_id"`
	Title       null.String `db:"title" json:"title"`
	Description null.String `db:"description" json:"description"`
	Contact     null.String `db:"contact" json:"contact"`
	SortOrder   uint32      `db:"sort_order" json:"sort_order"`
}

func (q *Queries) CreateTicket(ctx context.Context, arg CreateTicketParams) error {
	_, err := q.db.ExecContext(ctx, createTicket,
		arg.StatusID,
		arg.Title,
		arg.Description,
		arg.Contact,
		arg.SortOrder,
	)
	return err
}

const getLastInsertTicketByStatusID = `-- name: GetLastInsertTicketByStatusID :one
SELECT
  id, status_id, title, description, contact, sort_order, created_at, updated_at
FROM
  tickets
WHERE
  tickets.status_id = ?
  AND id = (
    SELECT
      LAST_INSERT_ID()
    FROM
      tickets AS t
    LIMIT
      1
  )
`

func (q *Queries) GetLastInsertTicketByStatusID(ctx context.Context, statusID uint32) (Ticket, error) {
	row := q.db.QueryRowContext(ctx, getLastInsertTicketByStatusID, statusID)
	var i Ticket
	err := row.Scan(
		&i.ID,
		&i.StatusID,
		&i.Title,
		&i.Description,
		&i.Contact,
		&i.SortOrder,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTicketByID = `-- name: GetTicketByID :one
SELECT
  id, status_id, title, description, contact, sort_order, created_at, updated_at
FROM
  tickets
WHERE
  id = ?
`

func (q *Queries) GetTicketByID(ctx context.Context, id uint64) (Ticket, error) {
	row := q.db.QueryRowContext(ctx, getTicketByID, id)
	var i Ticket
	err := row.Scan(
		&i.ID,
		&i.StatusID,
		&i.Title,
		&i.Description,
		&i.Contact,
		&i.SortOrder,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTicketWithBoard = `-- name: GetTicketWithBoard :one
SELECT
  tickets.id, tickets.status_id, tickets.title, tickets.description, tickets.contact, tickets.sort_order, tickets.created_at, tickets.updated_at,
  boards.id, boards.user_id, boards.title, boards.sort_order, boards.created_at, boards.updated_at
FROM
  tickets
  JOIN statuses ON tickets.status_id = statuses.id
  JOIN boards ON statuses.board_id = boards.id
WHERE
  tickets.id = ?
  AND statuses.board_id = ?
  AND boards.user_id = ?
`

type GetTicketWithBoardParams struct {
	ID      uint64 `db:"id" json:"id"`
	BoardID uint32 `db:"board_id" json:"board_id"`
	UserID  uint64 `db:"user_id" json:"user_id"`
}

type GetTicketWithBoardRow struct {
	Ticket Ticket `db:"ticket" json:"ticket"`
	Board  Board  `db:"board" json:"board"`
}

func (q *Queries) GetTicketWithBoard(ctx context.Context, arg GetTicketWithBoardParams) (GetTicketWithBoardRow, error) {
	row := q.db.QueryRowContext(ctx, getTicketWithBoard, arg.ID, arg.BoardID, arg.UserID)
	var i GetTicketWithBoardRow
	err := row.Scan(
		&i.Ticket.ID,
		&i.Ticket.StatusID,
		&i.Ticket.Title,
		&i.Ticket.Description,
		&i.Ticket.Contact,
		&i.Ticket.SortOrder,
		&i.Ticket.CreatedAt,
		&i.Ticket.UpdatedAt,
		&i.Board.ID,
		&i.Board.UserID,
		&i.Board.Title,
		&i.Board.SortOrder,
		&i.Board.CreatedAt,
		&i.Board.UpdatedAt,
	)
	return i, err
}

const getTickets = `-- name: GetTickets :many
SELECT
  id, status_id, title, description, contact, sort_order, created_at, updated_at
FROM
  tickets
WHERE
  (
    status_id = coalesce(/*SLICE:status_ids*/?, status_id)
    OR status_id IN (/*SLICE:status_ids*/?)
  )
  AND (
    id = coalesce(/*SLICE:exclude_ids*/?, id)
    OR id NOT IN (/*SLICE:exclude_ids*/?)
  )
ORDER BY
  status_id ASC,
  (
    CASE
      WHEN ? = 'asc' THEN sort_order
    END
  ) ASC,
  (
    CASE
      WHEN ? = 'desc' THEN sort_order
    END
  ) DESC
`

type GetTicketsParams struct {
	StatusIds          []uint32    `db:"status_ids" json:"status_ids"`
	ExcludeIds         []uint64    `db:"exclude_ids" json:"exclude_ids"`
	SortOrderDirection interface{} `db:"sort_order_direction" json:"sort_order_direction"`
}

func (q *Queries) GetTickets(ctx context.Context, arg GetTicketsParams) ([]Ticket, error) {
	query := getTickets
	var queryParams []interface{}
	if len(arg.StatusIds) > 0 {
		for _, v := range arg.StatusIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:status_ids*/?", strings.Repeat(",?", len(arg.StatusIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:status_ids*/?", "NULL", 1)
	}
	if len(arg.StatusIds) > 0 {
		for _, v := range arg.StatusIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:status_ids*/?", strings.Repeat(",?", len(arg.StatusIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:status_ids*/?", "NULL", 1)
	}
	if len(arg.ExcludeIds) > 0 {
		for _, v := range arg.ExcludeIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:exclude_ids*/?", strings.Repeat(",?", len(arg.ExcludeIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:exclude_ids*/?", "NULL", 1)
	}
	if len(arg.ExcludeIds) > 0 {
		for _, v := range arg.ExcludeIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:exclude_ids*/?", strings.Repeat(",?", len(arg.ExcludeIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:exclude_ids*/?", "NULL", 1)
	}
	queryParams = append(queryParams, arg.SortOrderDirection)
	queryParams = append(queryParams, arg.SortOrderDirection)
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Ticket{}
	for rows.Next() {
		var i Ticket
		if err := rows.Scan(
			&i.ID,
			&i.StatusID,
			&i.Title,
			&i.Description,
			&i.Contact,
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

const getTicketsByBoardID = `-- name: GetTicketsByBoardID :many
SELECT
  tickets.id, status_id, tickets.title, description, contact, tickets.sort_order, tickets.created_at, tickets.updated_at, statuses.id, board_id, statuses.title, statuses.sort_order, statuses.created_at, statuses.updated_at
FROM
  tickets
  JOIN statuses ON tickets.status_id = statuses.id
WHERE
  statuses.board_id = ?
`

type GetTicketsByBoardIDRow struct {
	ID          uint64      `db:"id" json:"id"`
	StatusID    uint32      `db:"status_id" json:"status_id"`
	Title       null.String `db:"title" json:"title"`
	Description null.String `db:"description" json:"description"`
	Contact     null.String `db:"contact" json:"contact"`
	SortOrder   uint32      `db:"sort_order" json:"sort_order"`
	CreatedAt   null.Time   `db:"created_at" json:"created_at"`
	UpdatedAt   null.Time   `db:"updated_at" json:"updated_at"`
	ID_2        uint32      `db:"id_2" json:"id_2"`
	BoardID     uint32      `db:"board_id" json:"board_id"`
	Title_2     null.String `db:"title_2" json:"title_2"`
	SortOrder_2 uint32      `db:"sort_order_2" json:"sort_order_2"`
	CreatedAt_2 null.Time   `db:"created_at_2" json:"created_at_2"`
	UpdatedAt_2 null.Time   `db:"updated_at_2" json:"updated_at_2"`
}

func (q *Queries) GetTicketsByBoardID(ctx context.Context, boardID uint32) ([]GetTicketsByBoardIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getTicketsByBoardID, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTicketsByBoardIDRow{}
	for rows.Next() {
		var i GetTicketsByBoardIDRow
		if err := rows.Scan(
			&i.ID,
			&i.StatusID,
			&i.Title,
			&i.Description,
			&i.Contact,
			&i.SortOrder,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.BoardID,
			&i.Title_2,
			&i.SortOrder_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
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

const getTicketsByStatusID = `-- name: GetTicketsByStatusID :many
SELECT
  id, status_id, title, description, contact, sort_order, created_at, updated_at
FROM
  tickets
WHERE
  status_id = ?
ORDER BY
  sort_order ASC
`

func (q *Queries) GetTicketsByStatusID(ctx context.Context, statusID uint32) ([]Ticket, error) {
	rows, err := q.db.QueryContext(ctx, getTicketsByStatusID, statusID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Ticket{}
	for rows.Next() {
		var i Ticket
		if err := rows.Scan(
			&i.ID,
			&i.StatusID,
			&i.Title,
			&i.Description,
			&i.Contact,
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

const getTicketsWithBoard = `-- name: GetTicketsWithBoard :many
SELECT
  tickets.id, tickets.status_id, tickets.title, tickets.description, tickets.contact, tickets.sort_order, tickets.created_at, tickets.updated_at,
  boards.id, boards.user_id, boards.title, boards.sort_order, boards.created_at, boards.updated_at
FROM
  tickets
  JOIN statuses ON tickets.status_id = statuses.id
  JOIN boards ON statuses.board_id = boards.id
WHERE
  tickets.id IN (/*SLICE:ids*/?)
  AND statuses.board_id = ?
  AND boards.user_id = ?
`

type GetTicketsWithBoardParams struct {
	Ids     []uint64 `db:"ids" json:"ids"`
	BoardID uint32   `db:"board_id" json:"board_id"`
	UserID  uint64   `db:"user_id" json:"user_id"`
}

type GetTicketsWithBoardRow struct {
	Ticket Ticket `db:"ticket" json:"ticket"`
	Board  Board  `db:"board" json:"board"`
}

func (q *Queries) GetTicketsWithBoard(ctx context.Context, arg GetTicketsWithBoardParams) ([]GetTicketsWithBoardRow, error) {
	query := getTicketsWithBoard
	var queryParams []interface{}
	if len(arg.Ids) > 0 {
		for _, v := range arg.Ids {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:ids*/?", strings.Repeat(",?", len(arg.Ids))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:ids*/?", "NULL", 1)
	}
	queryParams = append(queryParams, arg.BoardID)
	queryParams = append(queryParams, arg.UserID)
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTicketsWithBoardRow{}
	for rows.Next() {
		var i GetTicketsWithBoardRow
		if err := rows.Scan(
			&i.Ticket.ID,
			&i.Ticket.StatusID,
			&i.Ticket.Title,
			&i.Ticket.Description,
			&i.Ticket.Contact,
			&i.Ticket.SortOrder,
			&i.Ticket.CreatedAt,
			&i.Ticket.UpdatedAt,
			&i.Board.ID,
			&i.Board.UserID,
			&i.Board.Title,
			&i.Board.SortOrder,
			&i.Board.CreatedAt,
			&i.Board.UpdatedAt,
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

const getTicketsWithMinimumSortOrder = `-- name: GetTicketsWithMinimumSortOrder :many
SELECT
  id, status_id, title, description, contact, sort_order, created_at, updated_at
FROM
  tickets
WHERE
  status_id = ?
  AND sort_order >= ?
ORDER BY
  (
    CASE
      WHEN ? = 'asc' THEN sort_order
    END
  ) ASC,
  (
    CASE
      WHEN ? = 'desc' THEN sort_order
    END
  ) DESC
`

type GetTicketsWithMinimumSortOrderParams struct {
	StatusID           uint32      `db:"status_id" json:"status_id"`
	SortOrder          uint32      `db:"sort_order" json:"sort_order"`
	SortOrderDirection interface{} `db:"sort_order_direction" json:"sort_order_direction"`
}

func (q *Queries) GetTicketsWithMinimumSortOrder(ctx context.Context, arg GetTicketsWithMinimumSortOrderParams) ([]Ticket, error) {
	rows, err := q.db.QueryContext(ctx, getTicketsWithMinimumSortOrder,
		arg.StatusID,
		arg.SortOrder,
		arg.SortOrderDirection,
		arg.SortOrderDirection,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Ticket{}
	for rows.Next() {
		var i Ticket
		if err := rows.Scan(
			&i.ID,
			&i.StatusID,
			&i.Title,
			&i.Description,
			&i.Contact,
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

const updateTicket = `-- name: UpdateTicket :exec
UPDATE
  tickets
SET
  status_id = ?,
  title = ?,
  description = ?,
  contact = ?,
  sort_order = ?,
  updated_at = NOW()
WHERE
  id = ?
`

type UpdateTicketParams struct {
	StatusID    uint32      `db:"status_id" json:"status_id"`
	Title       null.String `db:"title" json:"title"`
	Description null.String `db:"description" json:"description"`
	Contact     null.String `db:"contact" json:"contact"`
	SortOrder   uint32      `db:"sort_order" json:"sort_order"`
	ID          uint64      `db:"id" json:"id"`
}

func (q *Queries) UpdateTicket(ctx context.Context, arg UpdateTicketParams) error {
	_, err := q.db.ExecContext(ctx, updateTicket,
		arg.StatusID,
		arg.Title,
		arg.Description,
		arg.Contact,
		arg.SortOrder,
		arg.ID,
	)
	return err
}

const updateTicketSortOrderAndStatusID = `-- name: UpdateTicketSortOrderAndStatusID :exec
UPDATE
  tickets
SET
  sort_order = ?,
  status_id = ?,
  updated_at = CASE
    WHEN ? <> tickets.status_id THEN NOW()
    ELSE updated_at
  END
WHERE
  id = ?
`

type UpdateTicketSortOrderAndStatusIDParams struct {
	SortOrder uint32 `db:"sort_order" json:"sort_order"`
	StatusID  uint32 `db:"status_id" json:"status_id"`
	ID        uint64 `db:"id" json:"id"`
}

func (q *Queries) UpdateTicketSortOrderAndStatusID(ctx context.Context, arg UpdateTicketSortOrderAndStatusIDParams) error {
	_, err := q.db.ExecContext(ctx, updateTicketSortOrderAndStatusID,
		arg.SortOrder,
		arg.StatusID,
		arg.StatusID,
		arg.ID,
	)
	return err
}

const updateTicketStatusID = `-- name: UpdateTicketStatusID :exec
UPDATE
  tickets
SET
  status_id = ?,
  updated_at = NOW()
WHERE
  id = ?
`

type UpdateTicketStatusIDParams struct {
	StatusID uint32 `db:"status_id" json:"status_id"`
	ID       uint64 `db:"id" json:"id"`
}

func (q *Queries) UpdateTicketStatusID(ctx context.Context, arg UpdateTicketStatusIDParams) error {
	_, err := q.db.ExecContext(ctx, updateTicketStatusID, arg.StatusID, arg.ID)
	return err
}
