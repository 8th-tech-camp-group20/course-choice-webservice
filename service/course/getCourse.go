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

var paramInvalidGetCourseResp = types.CreateCourseResponse{
	Code: types.ParamInvalid,
}

var courseNotExitGetCourseResp = types.CreateCourseResponse{
	Code: types.CourseNotExisted,
}

func GetCourse(c *gin.Context) {
	fmt.Println("hhh")
	var getReq types.GetCourseRequest
	if err := c.ShouldBindQuery(&getReq); err != nil {
		fmt.Println("error1")
		fmt.Println(getReq.CourseID)
		c.JSON(http.StatusOK, paramInvalidGetCourseResp)
		return
	}

	thisCourse := model.Course{}

	getRes := database.MysqlDB.Table("course").Where("id=?", getReq.CourseID).Find(&thisCourse)

	fmt.Println(getReq.CourseID)
	fmt.Println(thisCourse.ID)

	//查找失败
	if getRes.Error != nil || strconv.FormatUint(uint64(thisCourse.ID), 10) != getReq.CourseID {
		c.JSON(http.StatusOK, courseNotExitGetCourseResp)
		return
	}

	var returnRes = types.GetCourseResponse{
		Code: types.OK,
		Data: struct {
			CourseID  string
			Name      string
			TeacherID string
		}{CourseID: strconv.FormatUint(uint64(thisCourse.ID), 10), Name: thisCourse.Name, TeacherID: strconv.FormatUint(uint64(thisCourse.TeacherId), 10)},
	}
	fmt.Println("error3")
	c.JSON(http.StatusOK, returnRes)
}
