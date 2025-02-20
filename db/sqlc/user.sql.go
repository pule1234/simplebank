// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
) VALUES (
             $1, $2, $3, $4
         ) RETURNING username, hashed_password, full_name, email, password_changed_at, create_at
`

type CreateUserParams struct {
	Username       string `db:"username"`
	HashedPassword string `db:"hashed_password"`
	FullName       string `db:"full_name"`
	Email          string `db:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreateAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
select username, hashed_password, full_name, email, password_changed_at, create_at from users
where username = $1 limit 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreateAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
    hashed_password = COALESCE($1, hashed_password),
    password_changed_at = COALESCE($2, password_changed_at),
    full_name = COALESCE($3, full_name),
    email = COALESCE($4, email)
where
    username = $5
    RETURNING username, hashed_password, full_name, email, password_changed_at, create_at
`

type UpdateUserParams struct {
	HashedPassword    sql.NullString `db:"hashed_password"`
	PasswordChangedAt sql.NullTime   `db:"password_changed_at"`
	FullName          sql.NullString `db:"full_name"`
	Email             sql.NullString `db:"email"`
	Username          string         `db:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.FullName,
		arg.Email,
		arg.Username,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreateAt,
	)
	return i, err
}
