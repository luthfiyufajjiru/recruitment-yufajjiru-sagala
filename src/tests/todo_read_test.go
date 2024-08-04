package tests

import (
	"database/sql"
	"errors"
	"net/http"
	"sagala-todo/dependency"
	"sagala-todo/pkg/constants"
	customerror "sagala-todo/pkg/custom-error"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var returnedCols = []string{"id", "content", "status", "created_at", "created_by", "updated_at", "updated_by", "deleted_at"}

func TestGetTaskDetail(t *testing.T) {
	type expectation struct {
		id          string
		expectError error
	}

	const commonStr = "some string"

	cfg := dependency.InitConfiguration()

	usecase := dependency.InitTodoUsecaseMock(cfg)

	mockDb := usecase.Sql[constants.ConnSqlDefault].MockCtrl

	expectations := []expectation{
		{
			id: commonStr,
		},
		{
			id:          commonStr,
			expectError: errors.New(commonStr),
		},
		{
			id:          "",
			expectError: sql.ErrNoRows,
		},
	}

	query := "SELECT (.+) FROM tasks WHERE \\(deleted_at is null OR deleted_at = 0\\) AND id = (.+)"

	for _, expectation := range expectations {
		if expectation.expectError != nil {
			mockDb.ExpectQuery(query).WithArgs(expectation.id).WillReturnError(expectation.expectError)
		} else if expectation.expectError == nil {
			mockDb.ExpectQuery(query).WithArgs(expectation.id).WillReturnRows(sqlmock.NewRows(returnedCols).AddRow(commonStr, commonStr, commonStr, 1, nil, 1, nil, 0))
		}

		_, err := usecase.GetTask(expectation.id)

		if expectation.id != "" && expectation.expectError != nil {
			assert.NotNil(t, err)
		} else if expectation.id != "" {
			assert.Nil(t, err)
		} else if expectation.id == "" {
			var msgError *customerror.HttpError
			ok := errors.As(err, &msgError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusNotFound, msgError.StatusCode)
		}
	}

	err := mockDb.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
