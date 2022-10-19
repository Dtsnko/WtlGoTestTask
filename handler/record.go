package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Dtsnko/model"
	"github.com/Dtsnko/request"
	"github.com/Dtsnko/storage"
	taskstorage "github.com/Dtsnko/taskstorage"
	"github.com/gocarina/gocsv"
)

type RecordHandler struct {
}

var TaskStorage taskstorage.TaskStorage

/*
func GetCustomQueryResult(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var request requests.RequestCustomQuery
	var records []Record
	json.Unmarshal(reqBody, &request)
	query := request.Query
	db.Raw(query).Scan(&records)
	json.NewEncoder(w).Encode(records)
}
*/

func (recHandler *RecordHandler) UploadRecords(w http.ResponseWriter, r *http.Request) {
	database := storage.New()
	err := database.Open()
	if err != nil {
		return
	}
	//myLoger := logger.NewDatabaseLogger(*database)
	var inside func(uint) = func(taskId uint) {
		//myLoger.LogEvent(taskId, "Starting")
		reqBody, _ := io.ReadAll(r.Body)
		var request request.RequestUploadContacts
		var curStat taskstorage.Statistics

		if err := json.Unmarshal(reqBody, &request); err != nil {
			curStat.Errors = append(curStat.Errors, err.Error())
			TaskStorage.TaskEnd(taskId, curStat)
			return
		}
		//myLoger.LogEvent(taskId, "Unmarshaling")

		curStat = readRecordsFromUrl(request.Url, request.ClientId)
		//myLoger.LogEvent(taskId, "Reading records")
		TaskStorage.TaskEnd(taskId, curStat)
	}
	taskId := TaskStorage.TaskRun(inside)
	json.NewEncoder(w).Encode(fmt.Sprintf("Task id is %d", taskId))
	database.Close()
}
func (recHandler *RecordHandler) GetRecords(w http.ResponseWriter, r *http.Request) {
	database := storage.New()
	err := database.Open()
	if err != nil {
		return
	}
	reqBody, _ := io.ReadAll(r.Body)
	var request request.RequestGetContacts
	var result []model.Record
	query := "SELECT * FROM records WHERE 1=1 "

	if err := json.Unmarshal(reqBody, &request); err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	if request.ClientId != "" {
		query += fmt.Sprintf("AND client_id LIKE '%s' ", request.ClientId)
	}
	if request.ContactNumber != "" {
		query += fmt.Sprintf("AND number LIKE '%s' ", request.ContactNumber)
	}
	if request.ContactName != "" {
		query += fmt.Sprintf("AND name LIKE '%s' ", request.ContactName)
	}

	database.Db.Raw(query).Scan(&result)
	database.Close()
	json.NewEncoder(w).Encode(result)
}

//Block helping functions

func readRecordsFromUrl(url string, clientId string) taskstorage.Statistics {
	var newRecords []model.Record
	var currentStatistics taskstorage.Statistics

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

func updateRecords(newRecords []model.Record, currStat *taskstorage.Statistics) {
	database := storage.New()
	err := database.Open()
	if err != nil {
		return
	}
	for i, value := range newRecords {
		if value.Number == "" || value.ClientId == "" {
			currStat.Errors = append(currStat.Errors, "Unrecognizable record in CSV. Record number: "+strconv.Itoa(i+1))
			continue
		}

		if database.Record().IsExist(value) {
			if !value.Available {
				database.Db.Model(model.Record{}).Where("Number LIKE ? AND client_id LIKE ?", value.Number, value.ClientId).Delete(model.Record{})
				currStat.ContactsDeleted++
			} else {
				database.Db.Model(model.Record{}).Where("Number LIKE ? AND client_id LIKE ?", value.Number, value.ClientId).Update(model.Record{Name: value.Name})
				currStat.ContactsUpdated++
			}
		} else {
			if !value.Available {
				continue
			}
			database.Record().Create(value)
			currStat.ContactsCreated++
		}
	}
	database.Close()

}
