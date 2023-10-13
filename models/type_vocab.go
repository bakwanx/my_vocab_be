package models

import "time"

type TypeVocab struct {
	IdType    int       `gorm:"primarykey" json:"id_type" form:"id_type"`
	Type      string    `json:"type" form:"type"`
	CreatedAt time.Time `form:"created_at"`
	UpdatedAt time.Time `form:"updated_at"`
}
