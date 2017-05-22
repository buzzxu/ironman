package ironman

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type DbConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

func CreateDB(dbConfig *DbConfig) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	))
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "t_" + defaultTableName
	}
	if err != nil {
		log.Panic(fmt.Errorf("Failed to connect to log mysql: %s", err))
	}
	db.DB().SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Hour)
	db.DB().Ping()
	db.LogMode(true)
	return db, err
}
