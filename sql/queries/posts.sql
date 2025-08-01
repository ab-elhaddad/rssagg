-- name: CreatePost :one
INSERT INTO
    posts (
        title,
        url,
        description,
        published_at,
        feed_id
    )
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (url) DO NOTHING
RETURNING
    *;

-- name: CreatePosts :many
INSERT INTO
    posts (
        title,
        url,
        description,
        published_at,
        feed_id
    )
VALUES (
        unnest($1::text []),
        unnest($2::text []),
        unnest($3::text []),
        unnest($4::timestamp[]),
        unnest($5::uuid [])
    )
ON CONFLICT (url) DO NOTHING
RETURNING
    *;

-- name: GetPostsForUser :many
SELECT posts.*
FROM posts
    JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE
    feed_follows.user_id = $1
ORDER BY posts.published_at DESC;