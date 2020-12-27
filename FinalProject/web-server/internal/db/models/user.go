package models

import "gorm.io/gorm"

//User is a model representing a record in the table of Users
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null"`
	Password string `gorm:"type:varchar(256);not null"`
}
