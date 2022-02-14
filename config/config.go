package config

import "github.com/spf13/viper"

type MySQLConfig struct {
	User, Pass, Host, Port string
}

// InitConfig 初始化配置
// path 表示额外配置文件路径，如测试中指定配置文件路径
func InitConfig(path ...string) (*MySQLConfig, string, string) {
	//var debug = true
	//// 设置配置文件信息
	//if debug == false {
	//	viper.SetConfigName("conf")
	//} else {
	//	viper.SetConfigName("conf_local")
	//}
	viper.SetConfigName("conf")
	viper.SetConfigType("json")
	// 设置配置文件目录
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	for i := range path {
		viper.AddConfigPath(path[i])
	}
	// 读入配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}
	// 设置 MySQL 配置
	MySQLConf := &MySQLConfig{
		User: viper.GetString("mysql.user"),
		Pass: viper.GetString("mysql.pass"),
		Host: viper.GetString("mysql.host"),
		Port: viper.GetString("mysql.port"),
	}
	host, port := viper.GetString("host"), viper.GetString("port")
	return MySQLConf, host, port
}
