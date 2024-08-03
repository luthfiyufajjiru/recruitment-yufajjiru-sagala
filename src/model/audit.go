package model

import (
	"sagala-todo/pkg/nullable"
	"time"
)

type (
	Audit struct {
		// fields *By might not be used right now but prepared it as a standard
		CreatedAt nullable.NullInt64  `json:"created_at" db:"created_at"`
		CreatedBy nullable.NullString `json:"created_by" db:"created_by"`
		UpdatedAt nullable.NullInt64  `json:"updated_at" db:"updated_at"`
		UpdatedBy nullable.NullString `json:"updated_by" db:"updated_by"`
		DeletedAt nullable.NullInt64  `json:"deleted_at" db:"deleted_at"`
	}
)

func (a *Audit) InsertNow() {
	now := time.Now().Unix()
	a.CreatedAt.SetValue(now)
	a.UpdatedAt.SetValue(now)
}

func (a *Audit) UpdateNow() {
	now := time.Now().Unix()
	a.UpdatedAt.SetValue(now)
}
