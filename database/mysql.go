package database

import (
	"course-choice-webservice/config"
	"course-choice-webservice/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlDB *gorm.DB

func InitDB(dbconf *config.MySQLConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/course?charset=utf8mb4&parseTime=True&loc=Local",
		dbconf.User, dbconf.Pass, dbconf.Host, dbconf.Port)
	var err error
	MysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	err = MysqlDB.AutoMigrate(&model.Member{})
	if err != nil {
		panic(err.Error())
	}
}
