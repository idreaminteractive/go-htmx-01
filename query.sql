-- name: GetTodo :one
SELECT * FROM todo
WHERE id = ? LIMIT 1;

-- name: ListTodos :many
SELECT * FROM todo
ORDER BY id;

-- name: CreateTodo :one
INSERT INTO todo (
  description,
  user_id
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateTodo :exec
UPDATE todo
set description = ?
WHERE id = ?;

-- name: SetTodoDone :exec
update todo
set status = true 
where id = ?;

-- name: DeleteTodo :exec
DELETE FROM todo
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