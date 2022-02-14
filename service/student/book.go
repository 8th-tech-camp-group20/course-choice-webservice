package student

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var (
	paramInvalidBookCourseRequest = types.BookCourseResponse{Code: types.ParamInvalid}
)

func BookCourse(c *gin.Context) {
	var req types.BookCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := bookCourseService(&req)
	c.JSON(http.StatusOK, *resp)
}

func bookCourseService(req *types.BookCourseRequest) *types.BookCourseResponse {
	studentId, err := strconv.ParseUint(req.StudentID, 10, 64)
	if err != nil {
		return &paramInvalidBookCourseRequest
	}
	courseId := req.CourseID
	//先根据studentId查信息
	stu := model.Member{
		Model: gorm.Model{
			ID: uint(studentId),
		},
	}
	rows := database.MysqlDB.First(&stu, "id = ?", studentId).RowsAffected
	////查询的信息为空，返回学生不存在
	if rows != 0 {
		return &types.BookCourseResponse{
			Code: types.StudentNotExisted,
		}
	}
	////该studentId对应type不为student，返回无权限
	if stu.UserType != types.Student {
		return &types.BookCourseResponse{
			Code: types.PermDenied,
		}
	}
	//再根据courseId查信息
	thisCourse := model.Course{}
	rows = database.MysqlDB.Table("course").Where("id=?", courseId).Find(&thisCourse).RowsAffected
	//查询的信息为空，返回课程不存在
	if rows == 0 {
		return &types.BookCourseResponse{
			Code: types.CourseNotExisted,
		}
	} else {
		//该courseId对应cap_remain为0，返回课程已选满
		if thisCourse.CapRemain == 0 {
			return &types.BookCourseResponse{
				Code: types.CourseNotAvailable,
			}
		}
	}

	//studentId、courseId都合法，就选中成功
	//courseId对应的课程cap_remain - 1
	thisCourse.CapRemain--
	database.MysqlDB.Table("course").Where("id = ?", courseId).Update("cap_remain", thisCourse.CapRemain)
	bookSuccess(studentId, courseId)
	return &types.BookCourseResponse{
		Code: types.OK,
	}
}

func bookSuccess(studentId uint64, courseId string) {
	//查询这个学生id是否存在
	stuCou := model.StudentCourse{
		StudentId: studentId,
	}
	rows := database.MysqlDB.Table("book_course").Where("student_id = ?", studentId).Find(&stuCou).RowsAffected

	//存在
	if rows != 0 {
		courseIdList := stuCou.CourseIdList
		newCourseIdList := courseIdList + "," + courseId
		stuCou.CourseIdList = newCourseIdList
		database.MysqlDB.Table("book_course").Where("student_id = ?", studentId).Updates(stuCou)
	} else {
		//不存在
		stuCou = model.StudentCourse{
			StudentId:    studentId,
			CourseIdList: courseId,
		}
		database.MysqlDB.Table("book_course").Create(&stuCou)
	}
}
