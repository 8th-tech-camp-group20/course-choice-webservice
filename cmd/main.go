package main

import (
	"course-choice-webservice/config"
	"course-choice-webservice/database"
	"course-choice-webservice/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置项
	dbConf, host, port := config.InitConfig()
	// 设置数据库
	database.InitDB(dbConf)
	eng := gin.Default()
	service.RegisterRouter(eng)
	if err := eng.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		fmt.Println(err.Error())
	}
}
