package databases

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var databaseHandler *gorm.DB = nil

func ConnectMysql(username, password, address, database string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, address, database)

	dbf, error := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if error != nil {
		panic(error.Error())
	}

	databaseHandler = dbf
}

func GetDatabaseHandleInstance() *gorm.DB {
	if databaseHandler == nil {
		fmt.Printf("you should connect the database before invoke `GetDatabaseHandleInstance` method\n\n")
	}
	return databaseHandler
}
