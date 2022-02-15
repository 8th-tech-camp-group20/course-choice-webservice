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

var paramInvalidBindCourseResp = types.BindCourseResponse{
	Code: types.ParamInvalid,
}

var courseHasBindResp = types.BindCourseResponse{
	Code: types.CourseHasBound,
}

//TODO 目前使用teacherID=0标记绑定
func BindCourse(c *gin.Context) {
	fmt.Println("hhh")
	var bindReq types.BindCourseRequest
	if err := c.ShouldBindJSON(&bindReq); err != nil {
		fmt.Println("error1")
		c.JSON(http.StatusOK, paramInvalidBindCourseResp)
		return
	}

	if bindReq.TeacherID == "0" {
		c.JSON(http.StatusOK, paramInvalidBindCourseResp)
		return
	}

	hasTeacher := model.Course{}

	findRes := database.MysqlDB.Table("course").Where("id=?", bindReq.CourseID).Find(&hasTeacher)

	//课程不存在
	if findRes.Error != nil || strconv.FormatUint(uint64(hasTeacher.ID), 10) != bindReq.CourseID {
		c.JSON(http.StatusOK, courseNotExitGetCourseResp)
		return
	}

	//课程已经绑定
	if hasTeacher.TeacherId != 0 {
		c.JSON(http.StatusOK, courseHasBindResp)
		return
	}

	updateRes := database.MysqlDB.Table("course").Model(&model.Course{}).Where("id=?", bindReq.CourseID).Update("teacher_id", bindReq.TeacherID)

	//更新失败
	if updateRes.Error != nil {
		c.JSON(http.StatusOK, types.BindCourseResponse{Code: types.CourseNotExisted})
		return
	}

	fmt.Println("error3")
	c.JSON(http.StatusOK, types.BindCourseResponse{Code: types.OK})

}
