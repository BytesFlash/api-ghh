package data

import (
	"time"

	"gorm.io/gorm"
)

type Log struct {
	ID           string         `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt    time.Time      `gorm:"default:now()"`
	UpdatedAt    time.Time      `gorm:"default:now()" json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Module       string         `json:"module"`
	UserNickname string         `json:"user_nickname"`
	Sensor       string         `json:"sensor"`
	PersonDni    string         `json:"person_dni"`
	Detail       string         `json:"detail"`
	Params       string         `json:"-"`
}

func (db DB) CreateLog(log *Log) {
	_ = db.Create(&log)
}

func (db DB) GetlogByUser(name string) (*Log, error) {
	var log *Log
	result := db.Where("user_nickname = ?", name).First(&log)
	if result.Error != nil {
		return log, result.Error
	}
	return log, nil
}

func (db DB) GetlogByModule(name string) (*Log, error) {
	var log *Log
	result := db.Where("module = ?", name).First(&log)
	if result.Error != nil {
		return log, result.Error
	}
	return log, nil
}

func (db DB) GetlogBySensor(name string) (*Log, error) {
	var log *Log
	result := db.Where("sensor = ?", name).First(&log)
	if result.Error != nil {
		return log, result.Error
	}
	return log, nil
}

func (db DB) GetLogs(date string, startHour string, finishHour string) ([]*Log, error) {
	var log []*Log
	startDate := date + " " + startHour
	finishDate := date + " " + finishHour

	result := db.Where("created_at BETWEEN ? AND ?", startDate, finishDate).Order("created_at asc").Find(&log)
	if result.Error != nil {
		return log, result.Error
	}
	return log, nil
}

func (db DB) GetAllLog() (log []*Log, err error) {
	result := db.Find(&log)
	if result.Error != nil {
		return log, result.Error
	}
	return
}
