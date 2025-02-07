-- name: CreateYear :one
INSERT INTO years (
  title,
  owner,
  description,
  start_date,
  end_date
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetYear :one
SELECT * FROM years
WHERE id = $1 LIMIT 1;

-- name: GetYearForUpdate :one
SELECT * FROM years
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: ListYears :many
SELECT * FROM years
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: AddYearBalance :one
UPDATE years
SET balance = balance + sqlc.arg(amount), final_balance = final_balance + sqlc.arg(final_amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteYear :exec
DELETE FROM years WHERE id = $1;

-- name: UpdateYear :one
UPDATE years
SET title = $2, description = $3, start_date = $4, end_date = $5
WHERE id = $1
RETURNING *;