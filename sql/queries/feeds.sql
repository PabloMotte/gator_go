-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetAllFeeds :many
SELECT feeds.id, feeds.created_at, feeds.updated_at, feeds.name, feeds.url, feeds.user_id, users.name as user_name
FROM feeds 
JOIN users on feeds.user_id = users.id
;

-- name: GetFeed :one
SELECT feeds.id, feeds.created_at, feeds.updated_at, feeds.name, feeds.url, feeds.user_id, users.name as user_name
FROM feeds 
JOIN users on feeds.user_id = users.id
WHERE feeds.url in ($1)
;

-- name: GetFeedByID :one
SELECT feeds.id, feeds.created_at, feeds.updated_at, feeds.name, feeds.url, feeds.user_id, users.name as user_name
FROM feeds 
JOIN users on feeds.user_id = users.id
WHERE feeds.id in ($1)
;
