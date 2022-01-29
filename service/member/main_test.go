package member

import (
	"course-choice-webservice/config"
	"course-choice-webservice/database"
	"testing"
)

func TestMain(m *testing.M) {
	// 加载配置项
	dbConf, _, _ := config.InitConfig("../../config")
	// 设置数据库
	database.InitDB(dbConf)
	m.Run()
}
