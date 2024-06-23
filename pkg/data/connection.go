package data

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

type ConnectionAuth struct {
	Database string
	UserName string
	Password string
	Port     string
	Host     string
	SSL      string
	TimeZone string
	SSLCa    string
	SSLCert  string
	SSLKey   string
}

func (conAuth ConnectionAuth) Connect() (*DB, error) {
	fmt.Println("Intentando conectar a la base de datos...")
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s sslrootcert=%s sslcert=%s sslkey=%s",
			conAuth.Host,
			conAuth.UserName,
			conAuth.Password,
			conAuth.Database,
			conAuth.Port,
			conAuth.SSL,
			conAuth.TimeZone,
			conAuth.SSLCa,
			conAuth.SSLCert,
			conAuth.SSLKey,
		),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		return nil, err
	}
	fmt.Println("Conexión a la base de datos exitosa:", db)
	return &DB{db}, nil
}

func Migrate(db DB) error {
	logrus.Printf("clean data flag: %t")

	/* db.Migrator().AddColumn(User{}, "Dni") */
	db.SetupJoinTable(Role{}, "Modules", RoleModule{})
	/* db.Migrator().AutoMigrate(&Institution{}) */ /*revisar falla*/
	/* db.Migrator().AutoMigrate(&AutentiaPerson{})
	db.Migrator().AutoMigrate(&Institution{}) */
	db.AutoMigrate(&User{})
	db.createTableIfNotExists(Module{})
	db.createTableIfNotExists(Function{})
	db.createTableIfNotExists(Role{})
	db.createTableIfNotExists(RoleModule{})
	db.createTableIfNotExists(User{})
	db.createTableIfNotExists(Password{})
	db.createTableIfNotExists(Log{})

	return nil
}

func (db DB) createTableIfNotExists(model interface{}) {
	if !db.Migrator().HasTable(model) {
		db.Migrator().CreateTable(model)
	}
}
