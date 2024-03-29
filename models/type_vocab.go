package models

import (
	"time"
)

type TypeVocab struct {
	IdType      int       `gorm:"primarykey" form:"id_type"`
	Type        string    `json:"type" form:"type"`
	Description string    `json:"description" form:"description"`
	CreatedAt   time.Time `form:"created_at"`
	UpdatedAt   time.Time `form:"updated_at"`
}
