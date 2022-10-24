package handler

import (
	"github.com/Dtsnko/storage"
	taskstorage "github.com/Dtsnko/taskstorage"
	"github.com/gorilla/mux"
)

type RequestHandler struct {
	recordHand  *RecordHandler
	taskHand    *TaskHandler
	taskStorage *taskstorage.TaskStorage
}

// Creation of new Handler for the incoming requests
// param: *taskStorage.TaskStorage
// param: *mux.Router
// Every Handler has storage for executing tasks. From it's interface you can run, end tasks
func New(taskStorage *taskstorage.TaskStorage) *RequestHandler {
	strg := storage.New()
	return &RequestHandler{recordHand: &RecordHandler{TaskStorage: taskStorage, DataStorage: strg}, taskHand: &TaskHandler{TaskStorage: taskStorage, DataStorage: strg}, taskStorage: taskStorage}
}

func (hndlr *RequestHandler) HandleFunctions(router *mux.Router) *RequestHandler {
	router.HandleFunc("/upload", hndlr.Record().UploadRecords).Methods("POST")
	router.HandleFunc("/get_contacts", hndlr.Record().GetRecords).Methods("POST")
	router.HandleFunc("/record_custom_query", hndlr.Record().GetCustomQueryResult).Methods("POST")
	router.HandleFunc("/get_task_result/{id}", hndlr.Task().GetTaskResult).Methods("GET")
	return hndlr
}

// Interface to get results FROM RECORDS database
func (hndlr *RequestHandler) Record() *RecordHandler {
	if hndlr.taskHand != nil {
		return hndlr.recordHand
	} else {
		return &RecordHandler{TaskStorage: hndlr.taskStorage}
	}
}

// Interface to get results from tasks database
func (hndlr *RequestHandler) Task() *TaskHandler {
	if hndlr.taskHand != nil {
		return hndlr.taskHand
	} else {
		return &TaskHandler{TaskStorage: hndlr.taskStorage}
	}
}
