package data

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Password struct {
	CreatedAt time.Time      `gorm:"primaryKey"`
	UpdatedAt time.Time      `gorm:"default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Password  string         `json:"password"`
	UserID    string         `json:"-" gorm:"default:NULL"`
	User      *User          `json:"user"`
}

type Passwords []Password

func (db DB) PasswordExists(userPassword string, userID string) error {
	var passwords *Passwords
	find := errors.New("Not found")
	db.Where("user_id = ?", userID).Find(&passwords)
	for _, password := range *passwords {
		if checkPasswordHash(userPassword, password.Password) {
			find = nil
			break
		}
	}
	return find
}

func (db DB) GetLastPasswordDate(userID string) (Password, error) {
	var lastPassword Password
	result := db.Where("user_id = ?", userID).Last(&lastPassword)
	if result.Error != nil {
		return lastPassword, result.Error
	}
	return lastPassword, nil
}

func (db DB) CreatePassword(pass string, user *User) {
	pass, _ = HashPassword(pass)
	password := Password{Password: pass, UserID: user.ID}
	_ = db.Create(&password)
}
