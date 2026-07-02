package models

import (
	"time"
)

type Admin struct {
	Id uint `gorm:"primaryKey;autoIncrement"`

	Username string `gorm:"index;unique"`
	Password string

	CreatedAt time.Time `gorm:"index"`
}
