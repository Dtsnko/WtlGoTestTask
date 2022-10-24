package main

import (
	"io/ioutil"
	"log"
	"net/http"

	handler "github.com/Dtsnko/handler"
	"github.com/Dtsnko/storage"
	"github.com/Dtsnko/taskstorage"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {

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

func ConfigureHandler() *handler.RequestHandler {
	return handler.New(taskstorage.New())
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	hndlr := ConfigureHandler()
	hndlr.HandleFunctions(myRouter)
	myRouter.HandleFunc("/", getFile).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))

}

// Imitation of file with url
func getFile(w http.ResponseWriter, r *http.Request) {
	fileBytes, _ := ioutil.ReadFile("TestCSV.csv")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}
