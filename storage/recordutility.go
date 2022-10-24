package storage

import "github.com/Dtsnko/model"

type RecordUtility struct {
	storage *Storage
}

func (recordUtility *RecordUtility) Create(record model.Record) error {
	recordUtility.storage.Db.Model(model.Record{}).Create(record)
	return nil
}

func (recordUtility *RecordUtility) Delete(record model.Record) error {
	recordUtility.storage.Db.Model(model.Record{}).Where("Number LIKE ? AND client_id LIKE ?", record.Number, record.ClientId).Delete(model.Record{})
	return nil
}

func (recordUtility *RecordUtility) Update(record model.Record) error {
	recordUtility.storage.Db.Model(model.Record{}).Where("Number LIKE ? AND client_id LIKE ?", record.Number, record.ClientId).Update(model.Record{Name: record.Name})
	return nil
}

func (recordUtility *RecordUtility) IsExist(record model.Record) bool {
	var count int
	recordUtility.storage.Db.Model(model.Record{}).Where("client_id = ? AND Number = ?", record.ClientId, record.Number).Count(&count)
	return count > 0
}

// Add validation for real FROM RECORDS query
func (recordUtility *RecordUtility) RawQuery(query string) []model.Record {
	var result []model.Record
	recordUtility.storage.Db.Raw(query).Scan(&result)
	return result
}
