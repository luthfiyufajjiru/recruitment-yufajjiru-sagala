package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sagala-todo/pkg/constants"
	customerror "sagala-todo/pkg/custom-error"
	"sagala-todo/pkg/nullable"
	"sagala-todo/src/model"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (u *TodoUsecase) PostTask(payload model.TaskDTO) (taskId string, err error) {
	if !payload.Content.Valid || payload.Content.String == "" {
		err = &customerror.HttpError{
			Message:    "content could not be empty",
			StatusCode: http.StatusBadRequest,
		}
		return
	}

	if !payload.Status.Valid || statusMap[payload.Status.String] < 1 {
		payload.Status.SetValue("Waiting List") // default status
	}

	taskId = uuid.NewString() // uuid is not the best option for simple task to perform but since we are using sqlite, uuid is simpliest implementation right now
	payload.InsertNow()

	query := squirrel.Insert(tasksTable).Columns("id", "content", "status", "created_at", "updated_at", "deleted_at").Values(taskId, payload.Content.String, payload.Status.String, payload.CreatedAt.Int64, payload.UpdatedAt.Int64, 0)
	queryStr, args := query.MustSql()
	var res sql.Result
	res, err = u.Sql[constants.ConnSqlDefault].Db.Exec(queryStr, args...)
	n, _ := res.RowsAffected()
	if err != nil {
		return
	} else if n < 1 {
		err = &customerror.HttpError{
			Message:    "failed posting task",
			StatusCode: http.StatusInternalServerError,
		}
		return
	}
	return
}
func (u *TodoUsecase) UpdateTask(taskId string, payload model.TaskDTO) (err error) {
	var taskIdRet string
	if !payload.Status.Valid || statusMap[payload.Status.String] < 1 {
		err = &customerror.HttpError{
			Message:    "invalid status",
			StatusCode: http.StatusBadRequest,
		}
		return
	}

	if !payload.Content.Valid || payload.Content.String == "" {
		err = &customerror.HttpError{
			Message:    "content could not be empty",
			StatusCode: http.StatusBadRequest,
		}
		return
	}

	payload.UpdateNow()

	query := squirrel.Update(tasksTable).
		Set("content", payload.Content.String).
		Set("status", payload.Status.String).
		Set("updated_at", payload.UpdatedAt).
		Where("id = ?", taskId).
		Where("(deleted_at is null OR deleted_at = 0)").
		Suffix("returning id")

	queryStr, args := query.MustSql()

	err = u.Sql[constants.ConnSqlDefault].Db.Get(&taskIdRet, queryStr, args...)
	errNoRow := errors.Is(err, sql.ErrNoRows)
	if err != nil && !errNoRow {
		return
	} else if err != nil && errNoRow {
		err = &customerror.HttpError{
			Message:    fmt.Sprintf("there is no task with id of %s", taskId),
			StatusCode: http.StatusNotFound,
		}
		return
	}
	return
}
func (u *TodoUsecase) DeleteTask(taskId string, isHardDelete bool) (err error) {
	now := time.Now().Unix()

	// toggle enable/disable hard delete from env
	if u.Config[constants.HardDelete] != "True" {
		isHardDelete = false
	}

	switch isHardDelete {
	case true:
		query := squirrel.Delete(tasksTable).Where("id = ?", taskId)
		queryStr, args := query.MustSql()
		var res sql.Result
		res, err = u.Sql[constants.ConnSqlDefault].Db.Exec(queryStr, args...)
		n, _ := res.RowsAffected()
		errNoRow := errors.Is(err, sql.ErrNoRows)
		if err != nil && !errNoRow {
			return
		} else if err != nil && errNoRow || n < 1 {
			err = &customerror.HttpError{
				Message:    fmt.Sprintf("there is no task with id of %s", taskId),
				StatusCode: http.StatusNotFound,
			}
			return
		}
		return
	default:
		var taskIdRet nullable.NullString
		query := squirrel.Update(tasksTable).Set("deleted_at", now).Set("updated_at", now).Where("(deleted_at is null OR deleted_at = 0)").Where("id = ?", taskId).Suffix("returning id")
		queryStr, args := query.MustSql()

		err = u.Sql[constants.ConnSqlDefault].Db.Get(&taskIdRet, queryStr, args...)
		errNoRow := errors.Is(err, sql.ErrNoRows)
		if err != nil && !errNoRow {
			return
		} else if err != nil && errNoRow {
			err = &customerror.HttpError{
				Message:    fmt.Sprintf("there is no task with id of %s", taskId),
				StatusCode: http.StatusBadRequest,
			}
			return
		}
	}

	return
}
