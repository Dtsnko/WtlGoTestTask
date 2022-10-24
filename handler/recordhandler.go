package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Dtsnko/model"
	"github.com/Dtsnko/request"
	"github.com/Dtsnko/storage"
	taskstorage "github.com/Dtsnko/taskstorage"
	"github.com/gocarina/gocsv"
)

type RecordHandler struct {
	//Points to parent handler.Handler task storage
	TaskStorage *taskstorage.TaskStorage
	DataStorage *storage.Storage
}

// region: Methods for incoming requests
func (recHandler *RecordHandler) GetCustomQueryResult(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var request request.RequestCustomQuery
	if err := json.Unmarshal(reqBody, &request); err != nil {
		return
	}
	query := request.Query
	err := recHandler.DataStorage.Open()
	if err != nil {
		return
	}
	records := recHandler.DataStorage.Record().RawQuery(query)
	json.NewEncoder(w).Encode(records)
	defer recHandler.DataStorage.Close()
}

func (recHandler *RecordHandler) UploadRecords(w http.ResponseWriter, r *http.Request) {
	taskId := recHandler.TaskStorage.TaskRun(recHandler.uploadRecords, *r)
	json.NewEncoder(w).Encode(fmt.Sprintf("Task id is %d", taskId))
}

func (recHandler *RecordHandler) uploadRecords(taskId uint, r http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var request request.RequestUploadContacts
	taskStat := recHandler.TaskStorage.TaskGetStatistics(taskId)
	if err := json.Unmarshal(reqBody, &request); err != nil {
		taskStat.Errors = append(taskStat.Errors, err)
		return
	}
	records, err := recHandler.readRecordsFromUrl(request.Url, request.ClientId)
	if err != nil {
		taskStat.Errors = append(taskStat.Errors, err)
		return
	}
	recHandler.TaskStorage.TaskSetStatistics(taskId, recHandler.uploadRecordsToDB(records))
	recHandler.TaskStorage.TaskStop(taskId)
}
func (recHandler *RecordHandler) GetRecords(w http.ResponseWriter, r *http.Request) {
	err := recHandler.DataStorage.Open()
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

	recHandler.DataStorage.Db.Raw(query).Scan(&result)
	json.NewEncoder(w).Encode(result)
	defer recHandler.DataStorage.Close()
}

//region: Helping functions

func (recHandler *RecordHandler) readRecordsFromUrl(url string, clientId string) ([]model.Record, error) {
	var newRecords []model.Record

	response, httpError := http.Get(url)
	if httpError != nil {
		return nil, httpError
	}
	fileBytes, _ := io.ReadAll(response.Body)
	unmarshalError := gocsv.UnmarshalBytes(fileBytes, &newRecords)
	if unmarshalError != nil {
		return nil, unmarshalError
	}

	for index := range newRecords {
		newRecords[index].ClientId = clientId
	}
	return newRecords, nil
}

func (recHandler *RecordHandler) uploadRecordsToDB(newRecords []model.Record) taskstorage.Statistics {
	err := recHandler.DataStorage.Open()
	stats := taskstorage.Statistics{}
	if err != nil {
		return stats
	}
	for i, value := range newRecords {
		if value.Number == "" || value.ClientId == "" {
			stats.Errors = append(stats.Errors, errors.New("Unrecognizable record in CSV. Record number: "+strconv.Itoa(i+1)))
			continue
		}

		if recHandler.DataStorage.Record().IsExist(value) {
			if !value.Available {
				recHandler.DataStorage.Record().Delete(value)
				stats.Deleted++
			} else {
				recHandler.DataStorage.Record().Update(value)
				stats.Updated++
			}
		} else {
			if !value.Available {
				continue
			}
			recHandler.DataStorage.Record().Create(value)
			stats.Created++
		}
	}
	defer recHandler.DataStorage.Close()
	return stats
}
