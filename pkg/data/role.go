package data

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Role struct {
	ID        string         `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `gorm:"default:now()" json:"-"`
	UpdatedAt time.Time      `gorm:"default:now()" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null"`
	Modules   []*Module      `gorm:"many2many:role_module;"`
}

type RoleModule struct {
	RoleID    string         `gorm:"primaryKey; uniqueIndex:idx_role_module" `
	ModuleID  string         `gorm:"primaryKey; uniqueIndex:idx_role_module"`
	CreatedAt time.Time      `gorm:"default:now()"`
	UpdatedAt time.Time      `gorm:"default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (db DB) CreateRole(role *Role) {
	_ = db.Create(&role)
}

func (db DB) GetRole(name string) (*Role, error) {
	var role *Role
	result := db.Where("name = ?", name).First(&role)
	if result.Error != nil {
		return role, result.Error
	}
	return role, nil
}

func (db DB) GetAllRole() (role []*Role, err error) {
	result := db.Find(&role)
	if result.Error != nil {
		return role, result.Error
	}
	return
}

func (db DB) UpdateRole(role *Role, newName string) (*Role, error) {
	if newName != "" {
		_ = db.Model(&role).Where("name = ?", role.Name).Update("name", newName)
		role, _ := db.GetRole(newName)
		return role, nil
	} else {
		return &Role{}, errors.New("empty new name")
	}
}

func (db DB) GetRoleByName(name string) (role *Role, err error) {
	result := db.Where("name = ?", name).First(&role)
	if result.Error != nil {
		err = result.Error
	}
	return
}

func (db DB) GetRoleByID(id string) (role *Role, err error) {
	result := db.Where("id = ?", id).First(&role)
	if result.Error != nil {
		err = result.Error
	}
	return
}

func (db DB) CreateDefaultRole(name string) (*Role, error) {
	var role *Role
	role, err := db.GetRoleByName(name)
	if err != nil {
		if err.Error() == "record not found" {
			role = &Role{
				Name: name,
			}
			db.CreateRole(role)
			logrus.Printf("role: %s, created!", name)
			return role, nil
		} else {
			return nil, err
		}
	}

	if role.Name == name {
		logrus.Printf("role: %s, exists!", name)
	}
	return role, err
}

func (db DB) CreateRoleModule(roleModule *RoleModule) {
	_ = db.Create(&roleModule)
}

func (db DB) GetRoleModuleByIDs(idRol string, idModule string) (roleModule *RoleModule, err error) {
	result := db.Where("role_id = ?", idRol).Where("module_id = ?", idModule).First(&roleModule)
	if result.Error != nil {
		err = result.Error
	}
	return
}

func (db DB) CreateDefaultRoleModule(idRol string, idModule string) (*RoleModule, error) {
	role, _ := db.GetRoleByID(idRol)
	module, _ := db.GetModuleByID(idModule)
	var rm = "roleModule: [" + role.Name + ", " + module.Name + "], "
	var roleModule *RoleModule
	roleModule, err := db.GetRoleModuleByIDs(idRol, idModule)
	if err != nil {
		if err.Error() == "record not found" {
			roleModule = &RoleModule{
				RoleID:   idRol,
				ModuleID: idModule,
			}
			db.CreateRoleModule(roleModule)
			logrus.Printf("%s, created!", rm)
			return roleModule, nil
		} else {
			return nil, err
		}
	}

	if roleModule.RoleID == idRol && roleModule.ModuleID == idModule {
		logrus.Printf("%s, exists!", rm)
	}

	return roleModule, err
}

func (db DB) DeleteRoleModule(roleId string, moduleId string) error {
	var roleModule RoleModule
	db.Where("role_id = ? AND module_id = ? AND deleted_at IS NULL",
		roleId,
		moduleId,
	).First(&roleModule)
	result := db.Unscoped().Delete(&roleModule)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db DB) DeleteRole(roleId string) error {
	var role Role
	db.Where("id = ? AND deleted_at IS NULL", roleId).First(&role)
	result := db.Unscoped().Delete(&role)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
