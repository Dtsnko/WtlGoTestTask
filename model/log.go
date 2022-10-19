package model

import "github.com/jinzhu/gorm"

type Log struct {
	gorm.Model

	TaskId      uint
	Information string
}
