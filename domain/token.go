package domain

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Type string `gorm:"type:varchar(255)" json:"type"`
	Content string `gorm:"type:varchar(500)" json:"content"`
	Status string `gorm:"type:varchar(255)" json:"status"`
	UserID  int64 `json:"user_id"`  // Foreign key column
    User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // GORM relation
	CreatedAt time.Time `json:"created_at"`// auto set on insert
    UpdatedAt time.Time `json:"updated_at"` // auto set on update
}