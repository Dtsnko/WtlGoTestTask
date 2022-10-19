package taskstorage

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type task struct {
	Id       uint       `json:"id"`
	Status   string     `json:"status"`
	Function func(uint) `json:"-"`
	Stats    Statistics `json:"result"`
}

type TaskStorage struct {
	tasks []task
}

func New() *TaskStorage {
	return &TaskStorage{}
}

func (taskStorage *TaskStorage) TaskGetResult(taskId uint) Statistics {
	return taskStorage.tasks[taskId-1].Stats
}

func (taskStorage *TaskStorage) TaskRun(functionToRun func(taskId uint)) uint {
	id := uint(len(taskStorage.tasks) + 1)
	task := task{Id: id, Status: "Starting", Function: functionToRun}
	taskStorage.tasks = append(taskStorage.tasks, task)
	go task.Function(id)
	return id
}

func (taskStorage *TaskStorage) TaskEnd(taskId uint, stats Statistics) {
	taskStorage.tasks[taskId-1].Stats = stats
}

func (taskStorage *TaskStorage) GetTaskCount() int {
	return len(taskStorage.tasks)
}

/*
get Result

*/
