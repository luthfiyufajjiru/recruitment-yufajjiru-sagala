package usecase

import "sagala-todo/src/model"

func (u *TodoUsecase) PostTask(payload model.TaskDTO) (taskId string, err error) {
	return
}
func (u *TodoUsecase) UpdateTask(taskId string, payload model.TaskDTO) (err error) {
	return
}
func (u *TodoUsecase) DeleteTask(taskId string, isSoftDelete bool) (err error) {
	return
}
