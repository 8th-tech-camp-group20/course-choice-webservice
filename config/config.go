package config

import "github.com/spf13/viper"

type MySQLConfig struct {
	User, Pass, Host, Port string
}

func InitConfig() (*MySQLConfig, string, string) {
	// 设置配置文件信息
	viper.SetConfigName("conf")
	viper.SetConfigType("json")
	// 设置配置文件目录
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
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
