-- name: CreateUser :one
INSERT INTO
    users (
        name,
        email,
        password,
        api_key
    )
VALUES (
        $1,
        $2,
        $3,
        encode(
            sha256(random()::text::bytea),
            'hex'
        )
    )
RETURNING
    *;

-- name: LoginUser :one
SELECT * FROM users WHERE email = $1 AND password = $2;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;