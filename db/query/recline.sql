-- name: CreateRecLine :one
INSERT INTO reclines (
  title,
  owner,
  account_id,
  category_id,
  amount,
  description,
  recurrency,
  due_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetRecLine :one
SELECT * FROM reclines
WHERE id = $1 LIMIT 1;

-- name: GetRecLineForUpdate :one
SELECT * FROM reclines
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: ListRecLines :many
SELECT * FROM reclines
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteRecLine :exec
DELETE FROM reclines WHERE id = $1;

-- name: UpdateRecLine :one
UPDATE reclines
SET title = $2, account_id = $3, category_id = $4, amount = $5, description = $6, recurrency = $7, due_date = $8
WHERE id = $1
RETURNING *;
