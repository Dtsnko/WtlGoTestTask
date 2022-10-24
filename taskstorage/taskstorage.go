package taskstorage

import (
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type task struct {
	Id       uint                     `json:"id"`
	Status   string                   `json:"status"`
	Function func(uint, http.Request) `json:"-"`
	Stats    Statistics               `json:"result"`
}

type TaskStorage struct {
	tasks []task
}

func New() *TaskStorage {
	return &TaskStorage{}
}

func (taskStorage *TaskStorage) TaskGetStatistics(taskId uint) *Statistics {
	return &taskStorage.tasks[taskId-1].Stats
}
func (taskStorage *TaskStorage) TaskSetStatistics(taskId uint, stats Statistics) {
	taskStorage.tasks[taskId-1].Stats = stats
}

func (taskStorage *TaskStorage) TaskRun(functionToRun func(taskId uint, httpRequest http.Request), httpRequest http.Request) uint {
	id := uint(len(taskStorage.tasks) + 1)
	task := task{Id: id, Status: "Starting", Function: functionToRun, Stats: Statistics{}}
	taskStorage.tasks = append(taskStorage.tasks, task)
	go task.Function(id, httpRequest)
	return id
}

func (taskStorage *TaskStorage) TaskStop(taskId uint) {
	// todo: Functionality to stop tasks
}

func (taskStorage *TaskStorage) TaskGetCount() int {
	return len(taskStorage.tasks)
}
