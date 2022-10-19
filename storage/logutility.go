package storage

import "github.com/Dtsnko/model"

type LogUtility struct {
	storage *Storage
}

func (logUtility *LogUtility) Create(log model.Log) error {
	logUtility.storage.Db.Model(model.Log{}).Create(log)
	return nil
}
