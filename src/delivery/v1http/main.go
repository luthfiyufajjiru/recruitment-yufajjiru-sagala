package v1http

import (
	"sagala-todo/pkg/adapters"
	customlog "sagala-todo/pkg/custom-log"
	"sagala-todo/src/model"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger = customlog.Logger

type (
	Usecaser interface {
		PostTask(payload model.TaskDTO) (taskId string, err error)
		UpdateTask(taskId string, payload model.TaskDTO) (err error)
		DeleteTask(taskId string, isSoftDelete bool) (err error)
		GetTask(taskId string) (record model.TaskPresenter, err error)
		GetTasks() (records []model.TaskPresenter, err error)
	}

	V1Handler struct {
		Config  adapters.Config
		Usecase Usecaser
	}
)
