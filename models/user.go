package models

import "time"

type User struct {
	IdUser    int    `gorm:"primarykey" json:"id_user" form:"id_user"`
	Email     string `json:"email" form:"email"`
	Password  string `json:"password" form:"password"`
	Fullname  string `json:"fullname" form:"fullname"`
	Profile   string `json:"profile" form:"profile"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
