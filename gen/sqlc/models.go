// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package sqlc

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Author struct {
	ID   int64       `json:"id"`
	Name string      `json:"name"`
	Age  int32       `json:"age"`
	Bio  pgtype.Text `json:"bio"`
}