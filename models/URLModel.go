package models

import (
	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	ShortUrl   string `gorm:"unique"`
	MainUrl    string
	Fk_id_user int
	User       User `gorm:"foreignKey:Fk_id_user"`
}
