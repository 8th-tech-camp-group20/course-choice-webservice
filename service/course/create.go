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

var paramInvalidCreateCourseResp = types.CreateCourseResponse{
	Code: types.ParamInvalid,
}

func CreateCourse(c *gin.Context) {
	fmt.Println("hello")
	var courseReq types.CreateCourseRequest
	if err := c.ShouldBindJSON(&courseReq); err != nil {
		c.JSON(http.StatusBadRequest, paramInvalidCreateCourseResp)
		return
	}
	if l := len(courseReq.Name); l < 1 || courseReq.Cap < 0 {
		fmt.Println("error2")
		c.JSON(http.StatusBadRequest, paramInvalidCreateCourseResp)
		return
	}

	newCourse := model.Course{
		Name:      courseReq.Name,
		Cap:       courseReq.Cap,
		CapRemain: courseReq.Cap,
		TeacherId: 0,
	}

	//创建课程
	createRes := database.MysqlDB.Table("course").Create(&newCourse)

	//插入失败
	if createRes.Error != nil {
		c.JSON(http.StatusBadRequest, paramInvalidCreateCourseResp)
		return
	}

	var returnRes = types.CreateCourseResponse{
		Code: types.OK,
		Data: struct{ CourseID string }{CourseID: strconv.FormatUint(uint64(newCourse.ID), 10)},
	}
	fmt.Println("error3")
	c.JSON(http.StatusOK, returnRes)
}
