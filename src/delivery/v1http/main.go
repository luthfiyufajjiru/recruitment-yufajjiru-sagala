package v1http

import (
	"sagala-todo/pkg/adapters"
	"sagala-todo/src/model"
)

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
		Useacse Usecaser
	}
)
