package models

type ManagerModel struct {
	TaskModel TaskModelInterface
}

func NewManager() *ManagerModel {
	taskModel := newTaskModel()

	return &ManagerModel{
		TaskModel: taskModel,
	}
}