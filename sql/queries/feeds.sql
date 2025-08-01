-- name: CreateFeed :one
INSERT INTO
    feeds (name, url, user_id)
VALUES ($1, $2, $3)
RETURNING
    *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedsByUser :many
SELECT * FROM feeds WHERE user_id = $1;

-- name: GetFeedByID :one
SELECT * FROM feeds WHERE id = $1;

-- name: GetNextFeedsToFetch :many
SELECT id, url, name
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedAsFetched :exec
UPDATE feeds SET last_fetched_at = now() WHERE id = $1;