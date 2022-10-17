package main

import (
	"io/ioutil"
	"log"
	"net/http"

	record "github.com/Dtsnko/Record"
	task "github.com/Dtsnko/Task"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//Block routing

func main() {
	record.Db, _ = gorm.Open("sqlite3", "database.db")
	record.Db.AutoMigrate(&record.Record{})
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", getFile).Methods("GET")
	myRouter.HandleFunc("/upload", record.UploadContacts).Methods("POST")
	myRouter.HandleFunc("/get_contacts", record.GetContacts).Methods("POST")
	myRouter.HandleFunc("/custom_query", record.GetCustomQueryResult).Methods("POST")
	myRouter.HandleFunc("/get_task_result/{id}", task.TaskGetResult).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))

}

// Imitation of file with url
func getFile(w http.ResponseWriter, r *http.Request) {
	fileBytes, _ := ioutil.ReadFile("TestCSV.csv")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}
