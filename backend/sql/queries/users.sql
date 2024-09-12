-- name: CreateUser :one
INSERT INTO users (id, email, verified_email, name, given_name, family_name, picture, token, api_key)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?;

-- name: UpdateUser :one
UPDATE users
SET (verified_email, name, given_name, family_name, picture, token) = (?, ?, ?, ?, ?, ?)
WHERE email = ?
RETURNING *;

-- name: UpdateUserToken :one
UPDATE users
SET token = ?
WHERE email = ?
RETURNING *;

-- name: UpdateUserApiKey :one
UPDATE users
SET api_key = ?
WHERE email = ?
RETURNING *;

