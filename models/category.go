package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name  string
	Pid   int64
	Image string
}

func (*Category) TableName() string {
	return "category"
}
