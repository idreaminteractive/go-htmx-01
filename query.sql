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
WHERE id = ?
RETURNING *;

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

-- name: GetPublicNotes :many 
select n.id, n.content, u.id, n.updated_at, u.email
from note n, user u 
where 
  n.user_id = u.id and n.is_public = true;

-- test stuff

-- name: GetAllUsers :many
select * from user;


-- name: GetUserNoteAggregate :many
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
  10;