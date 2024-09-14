-- name: AddCategory :one
INSERT INTO categories (id, name)
VALUES (?, ?)
RETURNING *;

-- name: GetCategoryByName :one
SELECT *
FROM categories
WHERE name = ?;

-- name: GetCategories :many
SELECT *
FROM categories
ORDER BY name;
