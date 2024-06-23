package data

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Function struct {
	ID        string         `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `gorm:"default:now()"`
	UpdatedAt time.Time      `gorm:"default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"name" gorm:"not null"`
	ModuleID  string         `gorm:"index:idx_function_module" json:"-"`
	Module    *Module        `json:"module"`
	Active    bool           `gorm:"default:true;not null" json:"active"`
}

func (db DB) CreateFunction(function *Function) {
	_ = db.Create(&function)
}

func (db DB) GetFunctionByName(Name string) (function *Function, err error) {
	result := db.Where("name = ?", Name).First(&function)
	if result.Error != nil {
		err = result.Error
	}
	return
}

func (db DB) CreateDefaultFunction(name string, idModule string) (*Function, error) {
	var function *Function
	function, err := db.GetFunctionByName(name)
	if err != nil {
		if err.Error() == "record not found" {
			function = &Function{
				Name:     name,
				ModuleID: idModule,
			}
			db.CreateFunction(function)
			logrus.Printf("function: %s, created!", name)
			return function, nil
		} else {
			return nil, err
		}
	}

	if function.Name == name {
		logrus.Printf("function: %s, exists!", name)
	}

	return function, err
}
