package model

import "github.com/jinzhu/gorm"

type Record struct {
	gorm.Model
	//Id of client
	ClientId string

	//Contact properties
	Number    string `csv:"number"`
	Name      string `csv:"name"`
	Available bool   `csv:"available"`
}
