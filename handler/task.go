package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TaskHandler struct {
}

func (taskHandl *TaskHandler) GetTaskResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if id < 1 || id > TaskStorage.GetTaskCount()+1 {
		json.NewEncoder(w).Encode("Task doesn`t exist")
		return
	}
	json.NewEncoder(w).Encode(TaskStorage.TaskGetResult(uint(id)))
}
