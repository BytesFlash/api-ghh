package data

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Module struct {
	ID        string         `gorm:"primaryKey;default:uuid_generate_v4()" json:"-"`
	CreatedAt time.Time      `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time      `gorm:"default:now()" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null"`
	Functions []Function     `json:"functions"`
}

func (db DB) CreateModule(module *Module) {
	_ = db.Create(&module)
}

func (db DB) GetModuleByName(name string) (module *Module, err error) {
	result := db.Where("name = ?", name).First(&module)
	if result.Error != nil {
		err = result.Error
	}
	return
}

func (db DB) GetModuleByID(id string) (module *Module, err error) {
	result := db.Where("id = ?", id).First(&module)
	if result.Error != nil {
		err = result.Error
	}
	return
}

func (db DB) CreateDefaultModule(name string) (*Module, error) {
	var module *Module
	module, err := db.GetModuleByName(name)
	if err != nil {
		if err.Error() == "record not found" {
			module = &Module{
				Name: name,
			}
			db.CreateModule(module)
			logrus.Printf("module: %s, created!", name)
			return module, nil
		} else {
			return nil, err
		}
	}

	if module.Name == name {
		logrus.Printf("module: %s, exists!", name)
	}

	return module, err
}

func (db DB) ListAllModules() (modules []*Module, err error) {
	result := db.Find(&modules)
	if result.Error != nil {
		err = result.Error
	}
	return
}

func (db DB) DeleteModule(module *Module) error {
	result := db.Delete(&module)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
