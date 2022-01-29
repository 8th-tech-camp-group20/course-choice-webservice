package model

import (
	"course-choice-webservice/types"
	"gorm.io/gorm"
)

// Member 对应 members 表
type Member struct {
	gorm.Model
	// 昵称, 用户名
	Nickname, Username string
	// 密码
	Password []byte
	// 用户类型
	UserType types.UserType
}
