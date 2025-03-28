package models

import "time"

type Ping struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Message   string    `gorm:"type:varchar(255)" json:"message"`
	CreatedAt time.Time `json:"created_at"`
}