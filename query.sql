

-- name: CreateUser :one
insert into user (
  password, email, handle
) values (? , ?, ?) returning *;


-- name: GetUserByEmail :one
select * from user 
where email = ? limit 1; 


-- name: GetAllUsers :many
select * from user;
