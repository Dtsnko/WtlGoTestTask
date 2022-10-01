package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gocarina/gocsv"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Record struct {
	gorm.Model
	//Id of client
	ClientId string

	//Contact properties
	Number    string `csv:"number"`
	Name      string `csv:"name"`
	Available bool   `csv:"available"`
}

// Type for get contacts request
type RequestGetContacts struct {
	ContactNumber string `json:"contactNumber"`
	ContactName   string `json:"contactName"`
	ClientId      string `json:"clientId"`
}

// Type for upload contacts request
type RequestUploadContacts struct {
	Url      string `json:"url"`
	ClientId string `json:"clientId"`
}

// Type for statistics
type Statistics struct {
	ContactsCreated int
	ContactsUpdated int
	ContactsDeleted int
	Errors          []string
}

type RequestCustomQuery struct {
	Query string `json:"query"`
}

type Task struct {
	Id       uint       `json:"id"`
	Status   string     `json:"status"`
	Function func(uint) `json:"-"`
	Result   Statistics `json:"result"`
}

//Block routing

func main() {
	db, _ = gorm.Open("sqlite3", "database.db")
	db.AutoMigrate(&Record{})

	handleRequests()
}

var Records []Record
var Tasks []Task
var db *gorm.DB

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", getFile).Methods("GET")
	myRouter.HandleFunc("/upload", uploadContacts).Methods("POST")
	myRouter.HandleFunc("/get_contacts", getContacts).Methods("POST")
	myRouter.HandleFunc("/custom_query", getCustomQueryResult).Methods("POST")
	myRouter.HandleFunc("/get_task_result/{id}", taskGetResult).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))

}

// Imitation of file with url
func getFile(w http.ResponseWriter, r *http.Request) {
	fileBytes, _ := ioutil.ReadFile("TestCSV.csv")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

// Block tasks

func taskGetResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if id < 0 || id > len(Tasks) {
		json.NewEncoder(w).Encode("Task doesn`t exist")
		return
	}
	json.NewEncoder(w).Encode(Tasks[id-1])
}

func taskRun(functionToRun func(uint)) uint {
	id := uint(len(Tasks) + 1)
	var task Task = Task{Id: id, Status: "Starting", Function: functionToRun}
	Tasks = append(Tasks, task)
	go task.Function(id)
	return id
}

func taskLog(taskId uint, newStatus string) {
	Tasks[taskId-1].Status = newStatus
}

func endTask(taskId uint, stats Statistics) {
	Tasks[taskId-1].Result = stats
}

//Block urls

func getCustomQueryResult(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var request RequestCustomQuery
	var records []Record
	json.Unmarshal(reqBody, &request)
	query := request.Query
	db.Raw(query).Scan(&records)
	json.NewEncoder(w).Encode(records)
}

func uploadContacts(w http.ResponseWriter, r *http.Request) {
	var inside func(uint) = func(taskId uint) {
		taskLog(taskId, "Reading request")
		reqBody, _ := io.ReadAll(r.Body)
		var request RequestUploadContacts
		json.Unmarshal(reqBody, &request)
		taskLog(taskId, "Reading CSV file..")

		curStat := readRecordsFromUrl(request.Url, request.ClientId)
		taskLog(taskId, "Getting Stats")
		endTask(taskId, curStat)
		taskLog(taskId, "Ended")
	}
	taskId := taskRun(inside)
	json.NewEncoder(w).Encode(fmt.Sprintf("Task id is %d", taskId))
}
func getContacts(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var request RequestGetContacts
	var result []Record
	query := "SELECT * FROM records WHERE 1=1 "
	json.Unmarshal(reqBody, &request)
	if request.ClientId != "" {
		query += fmt.Sprintf("AND client_id LIKE '%s' ", request.ClientId)
	}
	if request.ContactNumber != "" {
		query += fmt.Sprintf("AND number LIKE '%s' ", request.ContactNumber)
	}
	if request.ContactName != "" {
		query += fmt.Sprintf("AND name LIKE '%s' ", request.ContactName)
	}
	db.Raw(query).Scan(&result)
	json.NewEncoder(w).Encode(result)
}

//Block helping functions

func readRecordsFromUrl(url string, clientId string) Statistics {
	var newRecords []Record
	var currentStatistics Statistics

	response, httpError := http.Get(url)
	if httpError != nil {
		currentStatistics.Errors = append(currentStatistics.Errors, "Http request error. Error: "+httpError.Error())
		return currentStatistics
	}
	fileBytes, _ := io.ReadAll(response.Body)
	unmarshalError := gocsv.UnmarshalBytes(fileBytes, &newRecords)
	if unmarshalError != nil {
		currentStatistics.Errors = append(currentStatistics.Errors, "Error while unparsing CSV. Error: "+unmarshalError.Error())
		log.Fatal(unmarshalError)
		return currentStatistics
	}

	for index := range newRecords {
		newRecords[index].ClientId = clientId
	}
	updateRecords(newRecords, &currentStatistics)
	return currentStatistics
}

func isRecordExist(record Record) bool {
	var count int
	db.Model(Record{}).Where("client_id = ? AND Number = ?", record.ClientId, record.Number).Count(&count)

	return count > 0
}

func updateRecords(newRecords []Record, currStat *Statistics) {

	for i, value := range newRecords {
		if value.Number == "" || value.ClientId == "" {
			currStat.Errors = append(currStat.Errors, "Unrecognizable record in CSV. Record number: "+strconv.Itoa(i+1))
			continue
		}

		if isRecordExist(value) {
			if !value.Available {
				db.Model(Record{}).Where("Number LIKE ? AND client_id LIKE ?", value.Number, value.ClientId).Delete(Record{})
				currStat.ContactsDeleted++
			} else {
				db.Model(Record{}).Where("Number LIKE ? AND client_id LIKE ?", value.Number, value.ClientId).Update(Record{Name: value.Name})
				currStat.ContactsUpdated++
			}
		} else {
			if !value.Available {
				continue
			}
			db.Create(&value)
			currStat.ContactsCreated++
		}
	}
}
