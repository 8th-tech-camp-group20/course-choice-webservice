package model

import (
	"course-choice-webservice/types"
	"gorm.io/gorm"
)

// Member 对应 members 表
type Member struct {
	gorm.Model
	// 用户 ID, 昵称, 用户名, 密码
	Nickname, Username, Password string
	// 用户类型
	UserType types.UserType
}
