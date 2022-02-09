package model

import "gorm.io/gorm"

type StudentCourse struct {
	gorm.Model
	StudentId    uint64 `gorm:"primaryKey"`
	CourseIdList string `gorm:"column:course_id_list"`
}

func (StudentCourse) TableName() string {
	return "book_course"
}
