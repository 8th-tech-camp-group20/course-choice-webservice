package student

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

/*
	type GetStudentCourseRequest struct {
		StudentID string
	}

	type GetStudentCourseResponse struct {
		Code ErrNo
		Data struct {
			CourseList []TCourse
		}
	}
*/

var paramInvalidGetStudentcourseResp = types.GetStudentCourseResponse{
	Code: types.ParamInvalid,
	Data: struct{ CourseList []types.TCourse }{CourseList: nil},
}
var studentNotExistedStudentCourseResp = types.GetStudentCourseResponse{
	Code: types.StudentNotExisted,
	Data: struct{ CourseList []types.TCourse }{CourseList: nil},
}

var studentHasNoCourseStudentCourseResp = types.GetStudentCourseResponse{
	Code: types.StudentHasNoCourse,
	Data: struct{ CourseList []types.TCourse }{CourseList: nil},
}

func GetStudentCourse(c *gin.Context) {
	var req types.GetStudentCourseRequest
	req.StudentID = c.Query("studentID")
	fmt.Println(req.StudentID)
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := getStudentCourseService(&req)
	c.JSON(http.StatusOK, *resp)
}

func getStudentCourseService(req *types.GetStudentCourseRequest) *types.GetStudentCourseResponse {
	studentId, err := strconv.ParseUint(req.StudentID, 10, 64)
	if err != nil {
		return &paramInvalidGetStudentcourseResp
	}

	member := model.Member{}
	//若学生ID不存在 返回错误类型11
	result1 := database.MysqlDB.Table("members").Where("id=?", studentId).First(&member)
	if result1.Error != nil {
		return &studentNotExistedStudentCourseResp
	}
	//查找学生的课程id
	studentCourse := model.StudentCourse{CourseIdList: ""}
	_ = database.MysqlDB.Table("book_course").Where("student_id = ?", studentId).Find(&studentCourse)
	//若学生没有课程 返回错误类型13
	if len(studentCourse.CourseIdList) == 0 {
		return &studentHasNoCourseStudentCourseResp
	}
	courseIdList := strings.Split(studentCourse.CourseIdList, ",")
	courseList := []types.TCourse{}
	for _, courseId := range courseIdList {
		thisCourse := model.Course{}
		course_id, _ := strconv.ParseInt(courseId, 10, 0)
		_ = database.MysqlDB.Table("course").Where("id=?", course_id).Find(&thisCourse)
		courseList = append(courseList, struct {
			CourseID  string
			Name      string
			TeacherID string
		}{CourseID: courseId, Name: thisCourse.Name, TeacherID: strconv.FormatInt(thisCourse.TeacherId, 10)})
	}
	return &types.GetStudentCourseResponse{
		Code: types.OK,
		Data: struct{ CourseList []types.TCourse }{CourseList: courseList},
	}
}
