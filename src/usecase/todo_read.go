package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sagala-todo/pkg/constants"
	customerror "sagala-todo/pkg/custom-error"
	"sagala-todo/src/model"
	"strings"

	"github.com/Masterminds/squirrel"
)

var statusMap = map[string]int8{
	"Waiting List": 1,
	"In Progress":  1,
	"Done":         1,
}

func (u *TodoUsecase) GetTask(taskId string) (record model.TaskPresenter, err error) {
	query := squirrel.Select("*").From(tasksTable).Where("(deleted_at is null OR deleted_at = 0)").Where("id = ?", taskId)
	queryStr, args := query.MustSql()

	err = u.Sql[constants.ConnSqlDefault].Db.Get(&record, queryStr, args...)
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
func (u *TodoUsecase) GetTasks(limit, offset *int, search, status *string) (records []model.TaskPresenter, totalData int, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%w: %v", constants.ErrPanic, r)
		}
	}()

	if status != nil {
		_, available := statusMap[*status]
		if !available {
			err = &customerror.HttpError{
				Message:    "invalid status",
				StatusCode: http.StatusBadRequest,
			}
			return
		}
	}

	queryBuilder := func(col ...string) squirrel.SelectBuilder {
		query := squirrel.Select(col...).From(tasksTable).Where("(deleted_at is null OR deleted_at = 0)")

		if search != nil && *search != "" {
			*search = fmt.Sprintf("%%%s%%", strings.ToLower(*search))
			query = query.Where("content like ?", search)
		}

		if status != nil {
			logger.Info("helllo")
			query = query.Where("status = ?", status)
		}
		return query.OrderBy("created_at, id")
	}

	query := queryBuilder("*")
	query = query.Limit(uint64(*limit)).Offset(uint64(*offset))
	queryStr, args := query.MustSql()

	err = u.Sql[constants.ConnSqlDefault].Db.Select(&records, queryStr, args...)
	errNoRow := errors.Is(err, sql.ErrNoRows)
	if err != nil && !errNoRow {
		return
	} else if err != nil && errNoRow {
		err = nil
		return
	}

	queryCount := queryBuilder("count(content)")
	queryCountStr, argsCount := queryCount.MustSql()
	err = u.Sql[constants.ConnSqlDefault].Db.Get(&totalData, queryCountStr, argsCount...)
	errNoRow = errors.Is(err, sql.ErrNoRows)
	if err != nil && !errNoRow {
		return
	} else if err != nil && errNoRow {
		err = nil
		return
	}
	return
}
