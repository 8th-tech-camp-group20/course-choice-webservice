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

func GetStudentCourse(c *gin.Context) {
	var req types.GetStudentCourseRequest
	resp := getStudentCourseService(&req)
	c.JSON(http.StatusOK, *resp)
}

func getStudentCourseService(req *types.GetStudentCourseRequest) *types.GetStudentCourseResponse {
	studentId, err := strconv.ParseUint(req.StudentID, 10, 0)
	//var studentId = 1
	fmt.Println(studentId)
	if err != nil {
		return &paramInvalidGetStudentcourseResp
	}

	// 从student_course表中找到对应的课程id

	//获得学生id对应的Student_course结构体
	//_ = database.MysqlDB.Table("student").Where("student_id=?", studentId).Find(&studentCourse)
	//var teacherCourses []model.Course
	//_ = database.MysqlDB.Table("course").Where("id=?", 1).Find(teacherCourses)
	//fmt.Println(teacherCourses)
	studentCourse := model.StudentCourse{}
	_ = database.MysqlDB.Table("book_course").Where("student_id = ?", studentId).Find(&studentCourse)
	//fmt.Println(studentCourse)
	courseIdList := strings.Split(studentCourse.CourseIdList, ",")
	fmt.Println(courseIdList)
	courseList := []types.TCourse{}

	//_ = database.MysqlDB.Table("course").Where("id = ?", 1).Find(&thisCourse)
	//fmt.Println(thisCourse)
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

	//_ = database.MysqlDB.Table()Find(studentCourse, "student_id = ?", studentId)
	//提取结构体的Course_List 去掉逗号

	//courseList := []types.TCourse{}
	//for courseId := range courseIdList {
	//	append(courseList, struct {
	//		CourseID  string
	//		Name      string
	//		TeacherID string
	//	}{CourseID: string(courseId), Name: "", TeacherID: ""})
	//}

}
