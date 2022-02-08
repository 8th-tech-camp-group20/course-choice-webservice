package course

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var paramInvalidUnBindCourseResp = types.UnbindCourseResponse{
	Code: types.ParamInvalid,
}

var courseNotBindResp = types.UnbindCourseResponse{
	Code: types.CourseNotBind,
}

func UnbindCourse(c *gin.Context) {
	fmt.Println("hhh")
	var unbindReq types.UnbindCourseRequest
	if err := c.ShouldBindJSON(&unbindReq); err != nil {
		fmt.Println("error1")
		c.JSON(http.StatusBadRequest, paramInvalidUnBindCourseResp)
		return
	}

	hasTeacher := model.Course{}

	findRes := database.MysqlDB.Table("course").Where("id=?", unbindReq.CourseID).Find(&hasTeacher)

	//课程不存在
	if findRes.Error != nil {
		c.JSON(http.StatusBadRequest, courseNotExitGetCourseResp)
		return
	}

	//课程未绑定
	if hasTeacher.TeacherId == 0 {
		c.JSON(http.StatusBadRequest, courseNotBindResp)
		return
	}

	updateRes := database.MysqlDB.Table("course").Model(&model.Course{}).Where("id=?", unbindReq.CourseID).Update("teacher_id", 0)

	//更新失败
	if updateRes.Error != nil {
		c.JSON(http.StatusBadRequest, types.BindCourseResponse{Code: types.CourseNotExisted})
		return
	}

	fmt.Println("error3")
	c.JSON(http.StatusOK, types.BindCourseResponse{Code: types.OK})

}
