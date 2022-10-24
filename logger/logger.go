package logger

import (
	"fmt"
	"time"

	"github.com/Dtsnko/model"
	"github.com/Dtsnko/storage"
)

type DatabaseLogger struct {
	storage *storage.Storage
}

type ConsoleLogger struct {
	storage *storage.Storage
}

func (logger *DatabaseLogger) LogEvent(taskId uint, report string) {

	log := model.Log{Information: report, TaskId: taskId}
	logger.storage.Log().Create(log)
}
func (logger *ConsoleLogger) LogEvent(report string) {

	fmt.Println("Time: ", time.Now(), "Info: ", report)
}

func NewDatabaseLogger(storage storage.Storage) *DatabaseLogger {
	return &DatabaseLogger{storage: &storage}
}

func NewConsoleLogger(storage storage.Storage) *ConsoleLogger {
	return &ConsoleLogger{storage: &storage}
}
