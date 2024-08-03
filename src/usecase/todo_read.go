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
	return
}
func (u *TodoUsecase) GetTasks(limit, offset *int, search, status *string) (records []model.TaskPresenter, totalData int, err error) {
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
		query := squirrel.Select(col...).From(tasksTable)

		if search != nil && *search != "" {
			*search = fmt.Sprintf("%%%s%%", strings.ToLower(*search))
			query = query.Where("content like ?", search)
		}

		if status != nil {
			logger.Info("helllo")
			query = query.Where("status = ?", status)
		}
		return query
	}

	query := queryBuilder("*")
	query = query.Limit(uint64(*limit)).Offset(uint64(*offset))
	queryStr, args := query.MustSql()

	err = u.Sql[constants.ConnSqlDefault].Db.Select(&records, queryStr, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}

	queryCount := queryBuilder("count(content)")
	queryCountStr, argsCount := queryCount.MustSql()
	err = u.Sql[constants.ConnSqlDefault].Db.Get(&totalData, queryCountStr, argsCount...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	return
}
