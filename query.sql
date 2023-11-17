

-- name: CreateUser :one
insert into user (
  password, email, name
) values (? , ?) returning *;


-- name: GetUserByEmail :one
select * from user 
where email = ? limit 1; 
