package models

import "time"

type Role struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      string `json:"code" gorm:"varchar(15);not null"`
	Name      string `json:"name" gorm:"varchar(20);not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
