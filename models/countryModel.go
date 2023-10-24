package models

import "gorm.io/gorm"

type Country struct {
	gorm.Model
	Country string `json:"country"`
}
