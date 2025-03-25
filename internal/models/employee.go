package models

import "time"

type Employee struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement:true"`
	FirstName string    `json:"firstName" gorm:"not null"`
	LastName  string    `json:"lastName" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
