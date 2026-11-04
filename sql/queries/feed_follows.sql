-- name: CreateFeedFollow :one
with inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT feed_follows.id, feed_follows.created_at, feed_follows.updated_at, feed_follows.user_id, users.name as user_name, feed_follows.feed_id, feeds.name as feed_name, feeds.Url as feed_url
FROM inserted_feed_follow AS feed_follows 
JOIN users on feed_follows.user_id = users.id
JOIN feeds on feed_follows.feed_id = feeds.id
;

-- name: GetAllFollows :many
SELECT feed_follows.id, feed_follows.created_at, feed_follows.updated_at, feed_follows.user_id, users.name as user_name, feed_follows.feed_id, feeds.name as feed_name, feeds.Url as feed_url
FROM feed_follows 
JOIN users on feed_follows.user_id = users.id
JOIN feeds on feed_follows.feed_id = feeds.id
;

-- name: Unfollow :exec
DELETE
FROM feed_follows
WHERE feed_id = $1 AND user_id = $2
;

-- name: GetUserFollows :many
SELECT feed_follows.id, feed_follows.created_at, feed_follows.updated_at, feed_follows.user_id, users.name as user_name, feed_follows.feed_id, feeds.name as feed_name, feeds.Url as feed_url
FROM feed_follows
JOIN users on feed_follows.user_id = users.id
JOIN feeds on feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id IN ($1)
;

-- name: GetUsersFollowing :many
SELECT feed_follows.id, feed_follows.created_at, feed_follows.updated_at, feed_follows.user_id, users.name as user_name, feed_follows.feed_id, feeds.name as feed_name, feeds.Url as feed_url
FROM feed_follows
JOIN users on feed_follows.user_id = users.id
JOIN feeds on feed_follows.feed_id = feeds.id
WHERE feed_follows.feed_id IN ($1)
;
