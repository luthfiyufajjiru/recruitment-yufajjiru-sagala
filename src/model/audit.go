package model

import "sagala-todo/pkg/nullable"

type (
	Audit struct {
		CreatedAt nullable.NullInt64  `json:"created_at" db:"created_at"`
		CreatedBy nullable.NullString `json:"created_by" db:"created_by"`
		UpdatedAt nullable.NullInt64  `json:"updated_at" db:"updated_at"`
		UpdatedBy nullable.NullString `json:"updated_by" db:"updated_by"`
		DeletedAt nullable.NullInt64  `json:"deleted_at" db:"deleted_at"`
	}
)
