package model

import (
	"gorm.io/gorm"
)

// Cource 对应 course 表
type Course struct {
	gorm.Model
	// 课程名称
	Name string
	// 课程容量
	Cap int
	// 课程剩余容量
	CapRemain int
	// 对应课程老师
	TeacherId int64
}
