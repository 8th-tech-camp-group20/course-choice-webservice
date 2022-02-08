package model

type StudentCourse struct {
	StudentId    uint64
	CourseIdList string `gorm:"column:course_id_list"`
}

func (StudentCourse) TableName() string {
	return "book_course"
}
