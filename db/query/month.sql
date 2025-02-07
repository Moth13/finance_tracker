-- name: CreateMonth :one
INSERT INTO months (
  title,
  owner,
  description,
  year_id,
  start_date,
  end_date
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetMonth :one
SELECT * FROM months
WHERE id = $1 LIMIT 1;

-- name: GetMonthForUpdate :one
SELECT * FROM months
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: ListMonths :many
SELECT * FROM months
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: AddMonthBalance :one
UPDATE months
SET balance = balance + sqlc.arg(amount), final_balance = final_balance + sqlc.arg(final_amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteMonth :exec
DELETE FROM months WHERE id = $1;

-- name: UpdateMonth :one
UPDATE months
SET title = $2, description = $3, year_id = $4, start_date = $5, end_date = $6
WHERE id = $1
RETURNING *;
