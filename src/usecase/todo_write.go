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
)

func (u *TodoUsecase) PostTask(payload model.TaskDTO) (taskId string, err error) {
	return
}
func (u *TodoUsecase) UpdateTask(taskId string, payload model.TaskDTO) (err error) {
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
