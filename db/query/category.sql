-- name: CreateCategory :one
INSERT INTO categories (
  title,
  owner
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: GetCategoryForUpdate :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: ListCategories :many
SELECT * FROM categories
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;