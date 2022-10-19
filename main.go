package main

import (
	"io/ioutil"
	"log"
	"net/http"

	handler "github.com/Dtsnko/handler"
	"github.com/Dtsnko/storage"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var hndlr *handler.Handler

//Block routing

func main() {
	hndlr = ConfigureHandler()

	handleRequests()
}

func ConfigureStorage() (*storage.Storage, error) {
	storage := storage.New()
	err := storage.Open()
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func ConfigureHandler() *handler.Handler {
	return handler.New()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", getFile).Methods("GET")
	myRouter.HandleFunc("/upload", hndlr.RecordHand.UploadRecords).Methods("POST")
	myRouter.HandleFunc("/get_contacts", hndlr.RecordHand.GetRecords).Methods("POST")
	//myRouter.HandleFunc("/custom_query", taskstorage.GetCustomQueryResult).Methods("POST")
	myRouter.HandleFunc("/get_task_result/{id}", hndlr.TaskHand.GetTaskResult).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))

}

// Imitation of file with url
func getFile(w http.ResponseWriter, r *http.Request) {
	fileBytes, _ := ioutil.ReadFile("TestCSV.csv")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}
