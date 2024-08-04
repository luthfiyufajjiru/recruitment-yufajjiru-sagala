package tests

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"net/http"
	"sagala-todo/dependency"
	"sagala-todo/pkg/constants"
	customerror "sagala-todo/pkg/custom-error"
	"sagala-todo/pkg/nullable"
	"sagala-todo/src/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/c2fo/testify/assert"
)

func TestPost(t *testing.T) {
	type expectation struct {
		payload          model.TaskDTO
		prepareStatement func(mdb sqlmock.Sqlmock)
		assertFn         func(t *testing.T, mdb sqlmock.Sqlmock, err error)
	}

	cfg := dependency.InitConfiguration()

	usecase := dependency.InitTodoUsecaseMock(cfg)

	mockDb := usecase.Sql[constants.ConnSqlDefault].MockCtrl

	expectations := []expectation{
		{
			payload: model.TaskDTO{},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				var msgErr *customerror.HttpError
				ok := errors.As(err, &msgErr)
				assert.True(t, ok)
				assert.Equal(t, http.StatusBadRequest, msgErr.StatusCode)
			},
		},
		{
			payload: model.TaskDTO{
				Content: nullable.NullString{
					NullString: sql.NullString{
						String: "hello",
						Valid:  true,
					},
				},
			},
			prepareStatement: func(mdb sqlmock.Sqlmock) {
				mdb.ExpectExec("INSERT INTO tasks (.+) VALUES (.+)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(driver.RowsAffected(1))
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.Nil(t, err)
				assert.Nil(t, mdb.ExpectationsWereMet())
			},
		},
		{
			payload: model.TaskDTO{
				Content: nullable.NullString{
					NullString: sql.NullString{
						String: "hello",
						Valid:  true,
					},
				},
			},
			prepareStatement: func(mdb sqlmock.Sqlmock) {
				mdb.ExpectExec("INSERT INTO tasks (.+) VALUES (.+)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(driver.RowsAffected(0))
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				var msgErr *customerror.HttpError
				ok := errors.As(err, &msgErr)
				assert.True(t, ok)
				assert.Equal(t, http.StatusInternalServerError, msgErr.StatusCode)
				assert.Nil(t, mdb.ExpectationsWereMet())
			},
		},
		{
			payload: model.TaskDTO{
				Content: nullable.NullString{
					NullString: sql.NullString{
						String: "hello",
						Valid:  true,
					},
				},
			},
			prepareStatement: func(mdb sqlmock.Sqlmock) {
				mdb.ExpectExec("INSERT INTO tasks (.+) VALUES (.+)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("some error"))
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, mdb.ExpectationsWereMet())
			},
		},
	}

	for _, expectation := range expectations {
		if expectation.payload.Content.Valid {
			expectation.prepareStatement(mockDb)
		}

		_, err := usecase.PostTask(expectation.payload)

		expectation.assertFn(t, mockDb, err)
	}
}

func TestUpdate(t *testing.T) {
	type expectation struct {
		payload          model.TaskDTO
		prepareStatement func(mdb sqlmock.Sqlmock)
		assertFn         func(t *testing.T, mdb sqlmock.Sqlmock, err error)
	}

	cfg := dependency.InitConfiguration()

	usecase := dependency.InitTodoUsecaseMock(cfg)

	mockDb := usecase.Sql[constants.ConnSqlDefault].MockCtrl

	query := "UPDATE tasks SET (.+) WHERE id = (.+) AND \\(deleted_at is null OR deleted_at = 0\\) returning id"

	expectations := []expectation{
		{
			payload: model.TaskDTO{},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				var msgErr *customerror.HttpError
				ok := errors.As(err, &msgErr)
				assert.True(t, ok)
				assert.Equal(t, http.StatusBadRequest, msgErr.StatusCode)
			},
		},
		{
			payload: model.TaskDTO{
				Content: nullable.NullString{
					NullString: sql.NullString{
						String: "hello",
						Valid:  true,
					},
				},
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				var msgErr *customerror.HttpError
				ok := errors.As(err, &msgErr)
				assert.True(t, ok)
				assert.Equal(t, http.StatusBadRequest, msgErr.StatusCode)
			},
		},
		{
			payload: model.TaskDTO{
				Content: nullable.NullString{
					NullString: sql.NullString{
						String: "hello",
						Valid:  true,
					},
				},
				Status: nullable.NullString{
					NullString: sql.NullString{
						String: "hello",
						Valid:  true,
					},
				},
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				var msgErr *customerror.HttpError
				ok := errors.As(err, &msgErr)
				assert.True(t, ok)
				assert.Equal(t, http.StatusBadRequest, msgErr.StatusCode)
			},
		},
		{
			payload: model.TaskDTO{
				Content: nullable.NullString{
					NullString: sql.NullString{
						String: "hello",
						Valid:  true,
					},
				},
				Status: nullable.NullString{
					NullString: sql.NullString{
						String: "Done",
						Valid:  true,
					},
				},
			},
			prepareStatement: func(mdb sqlmock.Sqlmock) {
				mdb.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("some_id"))
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.Nil(t, err)
				assert.Nil(t, mdb.ExpectationsWereMet())
			},
		},
		{
			payload: model.TaskDTO{
				Content: nullable.NullString{
					NullString: sql.NullString{
						String: "hello",
						Valid:  true,
					},
				},
				Status: nullable.NullString{
					NullString: sql.NullString{
						String: "Done",
						Valid:  true,
					},
				},
			},
			prepareStatement: func(mdb sqlmock.Sqlmock) {
				mdb.ExpectQuery(query).WillReturnError(errors.New("some error"))
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, mdb.ExpectationsWereMet())
			},
		},
		{
			payload: model.TaskDTO{
				Content: nullable.NullString{
					NullString: sql.NullString{
						String: "hello",
						Valid:  true,
					},
				},
				Status: nullable.NullString{
					NullString: sql.NullString{
						String: "Done",
						Valid:  true,
					},
				},
			},
			prepareStatement: func(mdb sqlmock.Sqlmock) {
				mdb.ExpectQuery(query).WillReturnError(sql.ErrNoRows)
			},
			assertFn: func(t *testing.T, mdb sqlmock.Sqlmock, err error) {
				var msgErr *customerror.HttpError
				ok := errors.As(err, &msgErr)
				assert.True(t, ok)
				assert.Equal(t, http.StatusNotFound, msgErr.StatusCode)
				assert.Nil(t, mdb.ExpectationsWereMet())
			},
		},
	}

	for _, expectation := range expectations {
		if expectation.prepareStatement != nil {
			expectation.prepareStatement(mockDb)
		}

		err := usecase.UpdateTask("some_id", expectation.payload)

		expectation.assertFn(t, mockDb, err)
	}
}
