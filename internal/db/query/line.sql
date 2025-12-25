-- name: CreateLine :one
INSERT INTO lines (
  title,
  owner,
  account_id,
  month_id,
  category_id,
  year_id,
  amount,
  checked,
  description,
  due_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetLine :one
SELECT * FROM lines
WHERE id = $1 LIMIT 1;

-- name: GetExpliciteLine :one
SELECT lines.id, lines.owner, lines.title, accounts.title as account, months.title as month, categories.title as category, lines.amount, lines.checked, lines.description, lines.due_date FROM lines
JOIN accounts ON accounts.id = lines.account_id
JOIN months ON months.id = lines.month_id
JOIN categories ON categories.id = lines.category_id
WHERE lines.id = $1;

-- name: GetLineForUpdate :one
SELECT * FROM lines
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: ListLines :many
SELECT * FROM lines
WHERE owner = $1
ORDER BY due_date DESC
LIMIT $2
OFFSET $3;

-- name: ListExplicitLines :many
SELECT lines.id, lines.owner, lines.title, accounts.title as account, months.title as month, categories.title as category, lines.amount, lines.checked, lines.description, lines.due_date FROM lines
JOIN accounts ON accounts.id = lines.account_id
JOIN months ON months.id = lines.month_id
JOIN categories ON categories.id = lines.category_id
WHERE lines.owner = $1
ORDER BY lines.due_date DESC
LIMIT $2
OFFSET $3;

-- name: DeleteLine :exec
DELETE FROM lines WHERE id = $1;

-- name: UpdateLine :one
UPDATE lines
SET title = $2, account_id = $3, month_id = $4, category_id = $5, year_id = $6, amount = $7, checked = $8, description = $9, due_date = $10
WHERE id = $1
RETURNING *;
