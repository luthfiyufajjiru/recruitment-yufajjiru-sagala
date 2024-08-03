package model

import "sagala-todo/pkg/nullable"

type (
	TaskDTO struct {
		Audit
		Content nullable.NullString `json:"content" db:"content"`
		Status  nullable.NullString `json:"status" db:"status"`
	}

	TaskPresenter struct {
		Audit
		Content nullable.NullString `json:"content" db:"content"`
		Status  nullable.NullString `json:"status" db:"status"`
		// This presentre migh be similar with DTO right now, but better to split due to additional response further development
	}
)
