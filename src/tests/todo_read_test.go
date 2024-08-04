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
	"github.com/c2fo/testify/assert"
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

func TestGetTasks(t *testing.T) {
	type expectation struct {
		limit, offset    *int
		search, status   *string
		prepareStatement func(mdb sqlmock.Sqlmock)
		assertFn         func(t *testing.T, mdb sqlmock.Sqlmock, err error)
	}

	cfg := dependency.InitConfiguration()

	usecase := dependency.InitTodoUsecaseMock(cfg)

	mockDb := usecase.Sql[constants.ConnSqlDefault].MockCtrl

	var (
		commonStr    = "somestr"
		commonStatus = "Done"
		commonInt    = 1
		commonError  = errors.New(commonStr)
	)

	expectations := []expectation{
		{
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.True(t, errors.Is(err, constants.ErrPanic))
			},
		},
		{
			status: &commonStr,
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				var msgError *customerror.HttpError
				ok := errors.As(err, &msgError)
				assert.True(t, ok)
				assert.Equal(t, http.StatusBadRequest, msgError.StatusCode)
			},
		},
		{
			status: &commonStatus,
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.True(t, errors.Is(err, constants.ErrPanic))
			},
		},
		{
			limit: &commonInt,
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.True(t, errors.Is(err, constants.ErrPanic))
			},
		},
		{
			offset: &commonInt,
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.True(t, errors.Is(err, constants.ErrPanic))
			},
		},
		{
			limit:  &commonInt,
			offset: &commonInt,
			prepareStatement: func(mdb sqlmock.Sqlmock) {
				mdb.ExpectQuery("SELECT (.+) FROM tasks WHERE (.+) ORDER BY (.+) LIMIT (.+) OFFSET (.+)").WillReturnError(commonError)
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.True(t, errors.Is(err, commonError))
				assert.Nil(t, mdb.ExpectationsWereMet())
			},
		},
		{
			limit:  &commonInt,
			offset: &commonInt,
			prepareStatement: func(mdb sqlmock.Sqlmock) {
				mdb.ExpectQuery("SELECT (.+) FROM tasks WHERE (.+) ORDER BY (.+) LIMIT (.+) OFFSET (.+)").WillReturnRows(sqlmock.NewRows(returnedCols).AddRow(commonStr, commonStr, commonStr, 1, nil, 1, nil, 0))
				mdb.ExpectQuery("SELECT (.+) FROM tasks WHERE (.+) ORDER BY (.+)").WillReturnError(commonError)
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.True(t, errors.Is(err, commonError))
				assert.Nil(t, mdb.ExpectationsWereMet())
			},
		},
		{
			limit:  &commonInt,
			offset: &commonInt,
			status: &commonStatus,
			prepareStatement: func(mdb sqlmock.Sqlmock) {
				mdb.ExpectQuery("SELECT (.+) FROM tasks WHERE (.+) AND (.+) ORDER BY (.+) LIMIT (.+) OFFSET (.+)").WillReturnRows(sqlmock.NewRows(returnedCols).AddRow(commonStr, commonStr, commonStr, 1, nil, 1, nil, 0))
				mdb.ExpectQuery("SELECT (.+) FROM tasks WHERE (.+) AND (.+) ORDER BY (.+)").WillReturnError(commonError)
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.True(t, errors.Is(err, commonError))
				assert.Nil(t, mdb.ExpectationsWereMet())
			},
		},
		{
			limit:  &commonInt,
			offset: &commonInt,
			status: &commonStatus,
			search: &commonStr,
			prepareStatement: func(mdb sqlmock.Sqlmock) {
				mdb.ExpectQuery("SELECT (.+) FROM tasks WHERE (.+) AND (.+) AND (.+) ORDER BY (.+) LIMIT (.+) OFFSET (.+)").WillReturnRows(sqlmock.NewRows(returnedCols).AddRow(commonStr, commonStr, commonStr, 1, nil, 1, nil, 0))
				mdb.ExpectQuery("SELECT (.+) FROM tasks WHERE (.+) AND (.+) AND (.+) ORDER BY (.+)").WillReturnError(commonError)
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.True(t, errors.Is(err, commonError))
				assert.Nil(t, mdb.ExpectationsWereMet())
			},
		},
		{
			limit:  &commonInt,
			offset: &commonInt,
			status: &commonStatus,
			search: &commonStr,
			prepareStatement: func(mdb sqlmock.Sqlmock) {
				mdb.ExpectQuery("SELECT (.+) FROM tasks WHERE (.+) AND (.+) AND (.+) ORDER BY (.+)").WillReturnError(commonError)
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.True(t, errors.Is(err, commonError))
				assert.Nil(t, mdb.ExpectationsWereMet())
			},
		},
		{
			limit:  &commonInt,
			offset: &commonInt,
			status: &commonStatus,
			search: &commonStr,
			prepareStatement: func(mdb sqlmock.Sqlmock) {
				mdb.ExpectQuery("SELECT (.+) FROM tasks WHERE (.+) AND (.+) AND (.+) ORDER BY (.+) LIMIT (.+) OFFSET (.+)").WillReturnRows(sqlmock.NewRows(returnedCols).AddRow(commonStr, commonStr, commonStr, 1, nil, 1, nil, 0))
				mdb.ExpectQuery("SELECT (.+) FROM tasks WHERE (.+) AND (.+) AND (.+) ORDER BY (.+)").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.Nil(t, err)
				assert.Nil(t, mdb.ExpectationsWereMet())
			},
		},
	}

	for _, expectation := range expectations {
		if expectation.limit != nil && expectation.offset != nil {
			expectation.prepareStatement(mockDb)
		}

		_, _, err := usecase.GetTasks(expectation.limit, expectation.offset, expectation.search, expectation.status)

		expectation.assertFn(t, mockDb, err)
	}

	err := mockDb.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}
