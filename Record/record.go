package record

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	requests "github.com/Dtsnko/Request"
	statistic "github.com/Dtsnko/Statistic"
	task "github.com/Dtsnko/Task"
	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
)

var Records []Record

var Db *gorm.DB

//Block urls

type Record struct {
	gorm.Model
	//Id of client
	ClientId string

	//Contact properties
	Number    string `csv:"number"`
	Name      string `csv:"name"`
	Available bool   `csv:"available"`
}

func GetCustomQueryResult(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var request requests.RequestCustomQuery
	var records []Record
	json.Unmarshal(reqBody, &request)
	query := request.Query
	Db.Raw(query).Scan(&records)
	json.NewEncoder(w).Encode(records)
}

func UploadContacts(w http.ResponseWriter, r *http.Request) {
	var inside func(uint) = func(taskId uint) {
		task.TaskLog(taskId, "Reading request")
		reqBody, _ := io.ReadAll(r.Body)
		var request requests.RequestUploadContacts
		json.Unmarshal(reqBody, &request)
		task.TaskLog(taskId, "Reading CSV file..")

		curStat := readRecordsFromUrl(request.Url, request.ClientId)
		task.TaskLog(taskId, "Getting Stats")
		task.EndTask(taskId, curStat)
		task.TaskLog(taskId, "Ended")
	}
	taskId := task.TaskRun(inside)
	json.NewEncoder(w).Encode(fmt.Sprintf("Task id is %d", taskId))
}
func GetContacts(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var request requests.RequestGetContacts
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
	Db.Raw(query).Scan(&result)
	json.NewEncoder(w).Encode(result)
}

//Block helping functions

func readRecordsFromUrl(url string, clientId string) statistic.Statistics {
	var newRecords []Record
	var currentStatistics statistic.Statistics

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
	Db.Model(Record{}).Where("client_id = ? AND Number = ?", record.ClientId, record.Number).Count(&count)

	return count > 0
}

func updateRecords(newRecords []Record, currStat *statistic.Statistics) {

	for i, value := range newRecords {
		if value.Number == "" || value.ClientId == "" {
			currStat.Errors = append(currStat.Errors, "Unrecognizable record in CSV. Record number: "+strconv.Itoa(i+1))
			continue
		}

		if isRecordExist(value) {
			if !value.Available {
				Db.Model(Record{}).Where("Number LIKE ? AND client_id LIKE ?", value.Number, value.ClientId).Delete(Record{})
				currStat.ContactsDeleted++
			} else {
				Db.Model(Record{}).Where("Number LIKE ? AND client_id LIKE ?", value.Number, value.ClientId).Update(Record{Name: value.Name})
				currStat.ContactsUpdated++
			}
		} else {
			if !value.Available {
				continue
			}
			Db.Create(&value)
			currStat.ContactsCreated++
		}
	}
}
