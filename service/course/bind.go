package course

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var paramInvalidBindCourseResp = types.BindCourseResponse{
	Code: types.ParamInvalid,
}

var courseHasBindResp = types.BindCourseResponse{
	Code: types.CourseHasBound,
}

func BindCourse(c *gin.Context) {
	fmt.Println("hhh")
	var bindReq types.BindCourseRequest
	if err := c.ShouldBindJSON(&bindReq); err != nil {
		fmt.Println("error1")
		c.JSON(http.StatusBadRequest, paramInvalidBindCourseResp)
		return
	}

	hasTeacher := model.Course{}

	findRes := database.MysqlDB.Table("course").Where("id=?", bindReq.CourseID).Find(&hasTeacher)

	//课程不存在
	if findRes.Error != nil {
		c.JSON(http.StatusBadRequest, courseNotExitGetCourseResp)
		return
	}

	//课程已经绑定
	if hasTeacher.TeacherId != 0 {
		c.JSON(http.StatusBadRequest, courseHasBindResp)
		return
	}

	updateRes := database.MysqlDB.Model(&model.Course{}).Where("id=?", bindReq.CourseID).Update("teacher_id", bindReq.TeacherID)

	//更新失败
	if updateRes.Error != nil {
		c.JSON(http.StatusBadRequest, types.BindCourseResponse{Code: types.CourseNotExisted})
		return
	}

	fmt.Println("error3")
	c.JSON(http.StatusOK, types.BindCourseResponse{Code: types.OK})

}
