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
		c.JSON(http.StatusOK, paramInvalidUnBindCourseResp)
		return
	}

	hasTeacher := model.Course{}

	findRes := database.MysqlDB.Table("course").Where("id=?", unbindReq.CourseID).Find(&hasTeacher)

	//课程不存在
	if findRes.Error != nil {
		c.JSON(http.StatusOK, courseNotExitGetCourseResp)
		return
	}

	//课程未绑定
	if hasTeacher.TeacherId == 0 || strconv.FormatUint(uint64(hasTeacher.TeacherId), 10) != unbindReq.TeacherID {
		c.JSON(http.StatusOK, courseNotBindResp)
		return
	}

	updateRes := database.MysqlDB.Table("course").Model(&model.Course{}).Where("id=?", unbindReq.CourseID).Update("teacher_id", 0)

	//更新失败
	if updateRes.Error != nil {
		c.JSON(http.StatusOK, types.BindCourseResponse{Code: types.CourseNotExisted})
		return
	}

	fmt.Println("error3")
	c.JSON(http.StatusOK, types.BindCourseResponse{Code: types.OK})

}
