// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package db

import (
	"context"
	"time"
)

const createNote = `-- name: CreateNote :one
INSERT INTO note (
  content,
  is_public,
  user_id
) VALUES (
  ?, ?, ?
)
RETURNING id, content, user_id, is_public, created_at, updated_at
`

type CreateNoteParams struct {
	Content  string
	IsPublic bool
	UserID   int64
}

func (q *Queries) CreateNote(ctx context.Context, arg CreateNoteParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, createNote, arg.Content, arg.IsPublic, arg.UserID)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
insert into user (
  password, email
) values (? , ?) returning id, first_name, last_name, password, email
`

type CreateUserParams struct {
	Password string
	Email    string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Password, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Password,
		&i.Email,
	)
	return i, err
}

const deleteNote = `-- name: DeleteNote :exec
DELETE FROM note
WHERE id = ?
`

func (q *Queries) DeleteNote(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteNote, id)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many

select id, first_name, last_name, password, email from user
`

// test stuff
func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Password,
			&i.Email,
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

const getNoteById = `-- name: GetNoteById :one
SELECT id, content, user_id, is_public, created_at, updated_at FROM note
WHERE id = ? LIMIT 1
`

func (q *Queries) GetNoteById(ctx context.Context, id int64) (Note, error) {
	row := q.db.QueryRowContext(ctx, getNoteById, id)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPublicNotes = `-- name: GetPublicNotes :many
select n.id, n.content, u.id, n.updated_at, u.email
from note n, user u 
where 
  n.user_id = u.id and n.is_public = true
`

type GetPublicNotesRow struct {
	ID        int64
	Content   string
	ID_2      int64
	UpdatedAt time.Time
	Email     string
}

func (q *Queries) GetPublicNotes(ctx context.Context) ([]GetPublicNotesRow, error) {
	rows, err := q.db.QueryContext(ctx, getPublicNotes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPublicNotesRow
	for rows.Next() {
		var i GetPublicNotesRow
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.ID_2,
			&i.UpdatedAt,
			&i.Email,
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

const getUserByEmail = `-- name: GetUserByEmail :one
select id, first_name, last_name, password, email from user 
where email = ? limit 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Password,
		&i.Email,
	)
	return i, err
}

const getUserNoteAggregate = `-- name: GetUserNoteAggregate :many
select
  user.id,
  user.email,
  json_group_array(json_object(
    'note_id', note.id,
    'content', note.content
   )) as notes,
   count(*) as num_notes
from
  user join note on note.user_id = user.id
  group by user.id
order by
  user.id
limit
  10
`

type GetUserNoteAggregateRow struct {
	ID       int64
	Email    string
	Notes    interface{}
	NumNotes int64
}

func (q *Queries) GetUserNoteAggregate(ctx context.Context) ([]GetUserNoteAggregateRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserNoteAggregate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserNoteAggregateRow
	for rows.Next() {
		var i GetUserNoteAggregateRow
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Notes,
			&i.NumNotes,
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

const listNotes = `-- name: ListNotes :many
SELECT id, content, user_id, is_public, created_at, updated_at FROM note
ORDER BY id
`

func (q *Queries) ListNotes(ctx context.Context) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, listNotes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.UserID,
			&i.IsPublic,
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

const listNotesForUser = `-- name: ListNotesForUser :many
SELECT (select count() from note) as count, id, content, user_id, is_public, created_at, updated_at FROM note
 where note.user_id = ? ORDER BY note.created_at desc limit 10
`

type ListNotesForUserRow struct {
	Count     int64
	ID        int64
	Content   string
	UserID    int64
	IsPublic  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) ListNotesForUser(ctx context.Context, userID int64) ([]ListNotesForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, listNotesForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListNotesForUserRow
	for rows.Next() {
		var i ListNotesForUserRow
		if err := rows.Scan(
			&i.Count,
			&i.ID,
			&i.Content,
			&i.UserID,
			&i.IsPublic,
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

const updateNote = `-- name: UpdateNote :exec
UPDATE note
set content = ?, is_public = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, content, user_id, is_public, created_at, updated_at
`

type UpdateNoteParams struct {
	Content  string
	IsPublic bool
	ID       int64
}

func (q *Queries) UpdateNote(ctx context.Context, arg UpdateNoteParams) error {
	_, err := q.db.ExecContext(ctx, updateNote, arg.Content, arg.IsPublic, arg.ID)
	return err
}
