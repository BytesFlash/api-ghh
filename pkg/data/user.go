package data

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/imedcl/manager-api/pkg/config"
)

type User struct {
	ID               string         `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt        time.Time      `gorm:"default:now()" json:"-"`
	UpdatedAt        time.Time      `gorm:"default:now()" json:"-"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	NickName         string         `gorm:"type:varchar(100);uniqueIndex;not null;default:uuid_generate_v4()" json:"nick_name"`
	Name             string         `json:"name"`
	Dni              string         `json:"dni"`
	Email            string         `json:"email"`
	Password         string         `json:"-"`
	Token            string         `gorm:"default:uuid_generate_v4()" json:"-"`
	ExpiresAt        time.Time      `gorm:"default:now()" json:"-"`
	Validated        bool           `json:"validated"`
	StatusUser       string         `gorm:"default:Activo" json:"status_user"`
	ActiInst         bool           `gorm:"default:false;not null" json:"acti_inst"`
	Active           bool           `gorm:"default:true;not null" json:"active"`
	Picture          string         `json:"picture"`
	Description      string         `json:"-"`
	LastTokenSession string         `json:"-"`
}

type Users []User

type Claims struct {
	User User `json:"user"`
	jwt.StandardClaims
}

const sessionDuration = 15

func (db DB) CreateUser(user *User) {
	user.Password, _ = HashPassword(user.Password)
	_ = db.Create(&user)
}

func (db DB) DeleteUser(user *User) error {
	result := db.Unscoped().Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db DB) Login(user *User, password string) (string, error) {
	cfg := config.New()
	if !user.Validated || user.StatusUser != config.STATUS_ACTIVE {
		return "", errors.New("Usuario inactivo")
	}

	mySigningKey := []byte(cfg.SignKey())
	if checkPasswordHash(password, user.Password) {
		//roles, _ := db.GetRoles(user.ID)
		//user.Roles = roles
		claims := &Claims{
			User: *user,
			StandardClaims: jwt.StandardClaims{
				Issuer: "Autentia Admin",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		var signed string
		signed, err := token.SignedString(mySigningKey)
		if err != nil {
			return "", err
		}
		user.ExpiresAt = time.Now().Add(time.Minute * sessionDuration)
		user.LastTokenSession = signed
		_ = db.Model(&user).Where("nick_name = ?", user.NickName).Updates(user)
		return signed, nil
	}
	return "", errors.New("Usuario - Contraseña no es una combinación válida")
}

func (db DB) ExtendToken(user *User) {
	user.ExpiresAt = time.Now().Add(time.Minute * sessionDuration)
	_ = db.Model(&user).Where("nick_name = ?", user.NickName).Updates(user)
}

func (db DB) GetRoleByUser(nickName string) (userRoleInstitution *User, err error) {
	result := db.Debug().Where("users.nick_name = ?", nickName).
		Preload("UserRolInstitutions.Institution").
		Preload("UserRolInstitutions.Role.Modules").
		Preload("Country").
		Find(&userRoleInstitution)
	if result.Error != nil {
		err = result.Error
	}

	return
}

func (db DB) UserExists(nickName string) bool {
	var user User
	result := db.Where("nick_name = ?", nickName).First(&user)
	return result.RowsAffected > 0
}

func (db DB) UserExistsEmail(email string) bool {
	var user User
	result := db.Where("email = ?", email).First(&user)
	return result.RowsAffected > 0
}
func (db DB) UserExistsName(nickName string) (user *User, err error) {
	result := db.Where("nick_name = ?", nickName).Preload("Country").First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return
}
func (db DB) UserExistEmail(email string) (user *User, err error) {
	result := db.Where("email = ?", email).Preload("Country").First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return
}

func (db DB) ExpireToken(user *User) {
	user.ExpiresAt = time.Now().Add(time.Minute * -1)
	_ = db.Model(&user).Where("nick_name = ?", user.NickName).Updates(user)
}

func (db DB) UserAllDBExists(nickName string) bool {
	var user User
	result := db.Unscoped().Where("nick_name = ?", nickName).First(&user)
	return result.RowsAffected > 0
}

func (db DB) GetUser(nickName string) (*User, error) {
	var user *User
	result := db.Where("nick_name = ?", nickName).Preload("Country").First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
func (db DB) GetUserEmail(email string) (*User, error) {
	var user *User
	result := db.Where("email = ?", email).Preload("Country").First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (db DB) GetUserReccovery(email string) (*User, error) {
	var user *User
	result := db.Where("email = ? ", email).Preload("Country").First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (db DB) GetUserByID(id string) (*User, error) {
	var user *User
	result := db.Where("ID = ?", id).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (db DB) IsActive(id string) bool {
	var user User
	result := db.Where("ID = ?", id).First(&user)
	if result.Error != nil {
		return false
	}
	return user.StatusUser == config.STATUS_ACTIVE && user.ExpiresAt.Unix() >= time.Now().Unix()
}

func (db DB) GetUserByToken(token string) (*User, error) {
	var user *User
	result := db.Where("Token = ?", token).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (db DB) UpdateUser(nickName string, user *User, password string) (*User, error) {
	if nickName != "" {
		if password != "" {
			user.Password, _ = HashPassword(password)
		}
		if user.Description == "" {
			_ = db.Model(&user).Where("nick_name = ?", nickName).Update("description", user.Description)
			_ = db.Model(&user).Where("nick_name = ?", nickName).Update("acti_inst", user.ActiInst)

		}
		_ = db.Model(&user).Where("nick_name = ?", nickName).Updates(user)
		if user.NickName != "" {
			return db.GetUser(user.NickName)
		} else {
			return db.GetUser(nickName)
		}

	} else {
		return &User{}, errors.New("user not found")
	}
}

func (db DB) GetUsers() *Users {
	var users *Users
	_ = db.Find(&users)

	return users
}

func (db DB) GetUsersByCountry(countryId string) []*User {
	var users []*User
	_ = db.Where("country_id = ?", countryId).
		Preload("Country").
		Find(&users)

	return users
}

func (db DB) GetUsersByRole(role string) []*User {
	admin, _ := db.GetRoleByName(role)
	var users []*User
	_ = db.
		Joins("JOIN user_role_institutions ON user_role_institutions.user_id = users.id").
		Where("user_role_institutions.role_id = ?", admin.ID).
		Preload("Country").
		Find(&users)

	return users
}

func (db DB) UserHasRole(role string, userId string) (user *User, err error) {
	admin, _ := db.GetRoleByName(role)
	result := db.Joins("JOIN user_role_institutions ON user_role_institutions.user_id = users.id").
		Where("user_role_institutions.role_id = ?", admin.ID).
		Where("users.id = ?", userId).
		Preload("Country").
		First(&user)

	if result.Error != nil {
		err = result.Error
	}
	return
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (db DB) CreateDefaultUser(nick string, name string, email string, pass string) (*User, error) {
	var user *User
	user, err := db.GetUser(nick)
	if err != nil {
		if err.Error() == "record not found" {
			user = &User{
				NickName:  nick,
				Name:      name,
				Email:     email,
				Password:  pass,
				Validated: true,
			}
			db.CreateUser(user)
			logrus.Printf("user: %s, created!", nick)
			return user, nil
		} else {
			return nil, err
		}
	}

	if user.NickName == nick {
		logrus.Printf("user: %s, exists!", nick)
	}

	return user, err
}
