-- name: AddPost :one
INSERT INTO posts (id, created_at, updated_at, 
            title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.id, posts.created_at, posts.updated_at, 
    posts.title, posts.Url, posts.description, posts.published_at,
    posts.feed_id, feeds.name as feed_name, 
    feed_follows.user_id, users.name as user_name
FROM posts
JOIN feed_follows on posts.feed_id = feed_follows.feed_id
JOIN feeds on feed_follows.feed_id = feeds.id
JOIN users on feed_follows.user_id = users.id
WHERE feed_follows.user_id IN ($1)
order by posts.created_at desc
OFFSET $2
LIMIT $3
;

-- name: GetPostByURL :one
SELECT posts.id, posts.created_at, posts.updated_at, 
    posts.title, posts.Url, posts.description, posts.published_at,
    posts.feed_id, feeds.name as feed_name
FROM posts 
JOIN feeds on posts.feed_id = feeds.id
WHERE posts.url in ($1)
LIMIT 1
;
