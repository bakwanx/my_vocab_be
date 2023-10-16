package models

import (
	"time"
)

// note
/*
foreignKey should name the model-local key field that joins to the foreign entity.
references should name the foreign entity's primary or unique key.
*/
type Vocab struct {
	IdVocab     int       `gorm:"primarykey" json:"id_vocab" form:"id_vocab"`
	IdUser      int       `json:"id_user" form:"id_user"`
	IdType      int       `json:"id_type" form:"id_type"`
	TypeVocab   TypeVocab `gorm:"foreignKey:IdType;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;references:IdType"`
	Vocab       string    `json:"vocab" form:"vocab"`
	Translation string    `json:"translation" form:"translation"`
	Variation   string    `json:"variation" form:"variation"`
	Note        string    `json:"note" form:"note"`
	CreatedAt   time.Time `form:"created_at"`
	UpdatedAt   time.Time `form:"updated_at"`
}
