package course

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var paramInvalidTeacherCourseResp = types.GetTeacherCourseResponse{
	Code: types.ParamInvalid,
}

var courseNotExitTeacherCourseResp = types.GetTeacherCourseResponse{
	Code: types.CourseNotExisted,
}

func GetTeacherCourse(c *gin.Context) {
	fmt.Println("hhh")
	var teacherReq types.GetTeacherCourseRequest
	if err := c.ShouldBindQuery(&teacherReq); err != nil {
		fmt.Println("error1")
		c.JSON(http.StatusOK, paramInvalidTeacherCourseResp)
		return
	}

	var teacherCourses []model.Course

	getTeacherRes := database.MysqlDB.Table("course").Where("teacher_id=?", teacherReq.TeacherID).Find(&teacherCourses)

	//查找失败
	if getTeacherRes.Error != nil {
		c.JSON(http.StatusOK, courseNotExitTeacherCourseResp)
		return
	}

	var returnRes []*types.TCourse

	for i := 0; i < len(teacherCourses); i++ {
		fmt.Println("isThis")
		var mid = types.TCourse{
			CourseID:  strconv.FormatUint(uint64(teacherCourses[i].ID), 10),
			Name:      teacherCourses[i].Name,
			TeacherID: strconv.FormatUint(uint64(teacherCourses[i].TeacherId), 10),
		}
		returnRes = append(returnRes, &mid)
	}

	fmt.Println("error3")
	c.JSON(http.StatusOK, types.GetTeacherCourseResponse{
		Code: types.OK,
		Data: struct{ CourseList []*types.TCourse }{CourseList: returnRes},
	})
}
