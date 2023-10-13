package models

import "time"

type Vocab struct {
	IdVocab   int       `gorm:"primarykey" json:"id_vocab" form:"id_vocab"`
	IdUser    int       `json:"id_user" form:"id_user"`
	IdType    string    `json:"id_type" form:"id_type"`
	Vocab     string    `json:"vocab" form:"vocab"`
	Variation string    `json:"variation" form:"variation"`
	Note      string    `json:"note" form:"note"`
	CreatedAt time.Time `form:"created_at"`
	UpdatedAt time.Time `form:"updated_at"`
}
