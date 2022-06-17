package database

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"omsoft.com/auth/cmd/models"
)

type DatabaseConfig struct {
	Driver   string
	Host     string
	Username string
	Password string
	Port     int
	Database string
}

type Database struct {
	*gorm.DB
}

func New(config *DatabaseConfig) (*Database, error) {
	var db *gorm.DB
	var err error
	switch strings.ToLower(config.Driver) {
	case "postgresql", "postgres":
		dsn := "user=" + config.Username + " password=" + config.Password + " dbname=" + config.Database + " host=" + config.Host + " port=" + strconv.Itoa(config.Port) + " TimeZone=UTC"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: true,
			CreateBatchSize: 100,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "ect_",                            // table name prefix, table for `User` would be `t_users`
				SingularTable: true,                              // use singular table name, table for `User` would be `user` with this option enabled
				NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
			},
			Logger: logger.Default.LogMode(logger.Warn)})
		if err != nil {
			panic(err)
		}
		err := db.AutoMigrate(&models.User{})
		if err != nil {
			fmt.Println("failed to automigrate User model:", err.Error())
			return &Database{db}, err
		}
	}
	return &Database{db}, err
}
