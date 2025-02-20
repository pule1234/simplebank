-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES (
             $1, $2, $3
         ) RETURNING *;

-- name: GetAccount :one
select * from accounts
where id  = $1 limit 1;

-- name: GetAccountForUpdate :one
select * from accounts
where id = $1 limit 1
for no key update;

-- name: ListAccounts :many
select * from accounts
where owner = $1
order by id
limit $2
offset $3;

-- name: UpdateAccount :one
update accounts
set balance = $2
where id = $1
RETURNING *;

-- name: AddAccountBalance :one
update accounts
set balance = balance + sqlc.arg(amount)
where id = sqlc.arg(id)
    RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;
