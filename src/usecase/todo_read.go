package usecase

import "sagala-todo/src/model"

func (u *TodoUsecase) GetTask(taskId string) (record model.TaskPresenter, err error) {
	return
}
func (u *TodoUsecase) GetTasks(limit, offset *int, search, status *string) (records []model.TaskPresenter, totalData int, err error) {
	return
}
