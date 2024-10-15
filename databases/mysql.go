package databases

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var databaseHandler *gorm.DB = nil

func ConnectMysql(username, password, database string) {
	dsn := fmt.Sprintf("%s:%s@/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, database)

	dbf, error := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if error != nil {
		panic(error.Error())
	}

	databaseHandler = dbf

	fmt.Printf("connect mysql %v\n", dbf)
}

func GetDatabaseHandleInstance() *gorm.DB {
	if databaseHandler == nil {
		fmt.Printf("you should connect the database before invoke `GetDatabaseHandleInstance` method\n\n")
	}
	return databaseHandler
}
