package storage

import "github.com/jinzhu/gorm"

type Storage struct {
	Db            *gorm.DB
	recordUtility *RecordUtility
	logUtility    *LogUtility
}

func New() *Storage {
	return &Storage{}
}

func (storage *Storage) Open() error {
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		return err
	}

	storage.Db = db
	return nil
}

func (storage *Storage) Close() {
	storage.Db.Close()
}

func (storage *Storage) Record() *RecordUtility {
	if storage.recordUtility != nil {
		return storage.recordUtility
	}
	return &RecordUtility{storage: storage}
}

func (storage *Storage) Log() *LogUtility {
	if storage.logUtility != nil {
		return storage.logUtility
	}
	return &LogUtility{storage: storage}
}
