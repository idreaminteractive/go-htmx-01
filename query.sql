-- name: GetNoteById :one
SELECT * FROM note
WHERE id = ? LIMIT 1;

-- name: ListNotes :many
SELECT * FROM note
ORDER BY id;


-- name: ListNotesForUser :many 
SELECT (select count() from note) as count, * FROM note
 where note.user_id = ? ORDER BY note.created_at desc limit 10;

-- name: CreateNote :one
INSERT INTO note (
  content,
  is_public,
  user_id
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: UpdateNote :exec
UPDATE note
set content = ?, is_public = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteNote :exec
DELETE FROM note
WHERE id = ?;


-- name: CreateUser :one
insert into user (
  password, email
) values (? , ?) returning *;


-- name: GetUserByEmail :one
select * from user 
where email = ? limit 1; 



-- test stuff

-- name: GetAllUsers :many
select * from user;