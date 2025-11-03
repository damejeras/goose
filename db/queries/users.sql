-- name: GetUser :one
select * from users where id = ?;

-- name: FindUserByGoogleID :one
select * from users where google_id = ?;

-- name: UpdateUserLastSeen :exec
update users set last_login_at = CURRENT_TIMESTAMP where id = ?;

-- name: UpdateUserProfile :exec
update users set name = ?, last_login_at = CURRENT_TIMESTAMP where id = ?;

-- name: CreateUser :one
insert into users (email, google_id, name) values (?, ?, ?) returning *;

