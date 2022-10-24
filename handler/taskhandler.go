package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Dtsnko/storage"
	"github.com/Dtsnko/taskstorage"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	//Points to parent handler.Handler task storage
	TaskStorage *taskstorage.TaskStorage
	//Points to storage where operations will execute
	DataStorage *storage.Storage
}

// region: Methods for incoming requests
func (taskHandlr *TaskHandler) GetTaskResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if id < 1 || id > taskHandlr.TaskStorage.TaskGetCount()+1 {
		json.NewEncoder(w).Encode("Task doesn`t exist")
		return
	}
	json.NewEncoder(w).Encode(taskHandlr.TaskStorage.TaskGetStatistics(uint(id)))
}
