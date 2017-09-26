package libs

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"log"
	"os"
)

var DB *gorm.DB

type DbConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func (c *DbConfig) InitDB() *gorm.DB {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
	db, err := gorm.Open("mysql", connString)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}

	return db
}
