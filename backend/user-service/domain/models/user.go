package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID        uuid.UUID `json:"uuid" gorm:"type:uuid;not null"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Username    string    `json:"username" gorm:"type:varchar(20);not null"`
	Password    string    `json:"password" gorm:"type:varchar(255);not null"`
	PhoneNumber string    `json:"phoneNumber" gorm:"type:varchar(15);not null"`
	Email       string    `json:"email" gorm:"type:varchar(100);not null"`
	RoleID      uint      `json:"roleId" gorm:"type:uint;not null"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time

	Role Role `json:"role" gorm:"foreignKey:role_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
