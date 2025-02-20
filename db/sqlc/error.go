package db

import (
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgconn"
)

// 错误码定义  postgreSql错误处理
const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

// 未查询到错误
var ErrRecordNotFound = pgx.ErrNoRows

var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
