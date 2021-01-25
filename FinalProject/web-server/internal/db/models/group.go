package models

import "time"

//User is a model representing a record in the table of Users
type Group struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(256);not null"`
	OwnerID   uint   `gorm:"type:Integer;not null"`
}
