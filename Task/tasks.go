package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	statistic "github.com/Dtsnko/Statistic"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var Tasks []Task

type Task struct {
	Id       uint                 `json:"id"`
	Status   string               `json:"status"`
	Function func(uint)           `json:"-"`
	Result   statistic.Statistics `json:"result"`
}

func TaskGetResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if id < 0 || id > len(Tasks) {
		json.NewEncoder(w).Encode("Task doesn`t exist")
		return
	}
	json.NewEncoder(w).Encode(Tasks[id-1])
}

func TaskRun(functionToRun func(taskId uint)) uint {
	id := uint(len(Tasks) + 1)
	var task Task = Task{Id: id, Status: "Starting", Function: functionToRun}
	Tasks = append(Tasks, task)
	go task.Function(id)
	return id
}

func TaskLog(taskId uint, newStatus string) {
	Tasks[taskId-1].Status = newStatus
}

func EndTask(taskId uint, stats statistic.Statistics) {
	Tasks[taskId-1].Result = stats
}
