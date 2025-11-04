-- name: CreateAPIKey :one
insert into api_keys (id, user_id, name, key_hash, key_prefix, key_suffix, created_at)
values (?, ?, ?, ?, ?, ?, current_timestamp)
returning *;

-- name: GetAPIKeyByHash :one
select * from api_keys
where key_hash = ?;

-- name: ListAPIKeysByUserID :many
select * from api_keys
where user_id = ?
order by created_at desc;

-- name: GetAPIKeyByID :one
select * from api_keys
where id = ? and user_id = ?;

-- name: DeleteAPIKey :exec
delete from api_keys
where id = ? and user_id = ?;

-- name: UpdateAPIKeyName :one
update api_keys
set name = ?
where id = ? and user_id = ?
returning *;

-- name: UpdateAPIKeyLastUsed :exec
update api_keys
set last_used_at = current_timestamp
where id = ?;
