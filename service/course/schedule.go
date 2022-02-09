package course

import (
	"course-choice-webservice/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var paramInvalidScheduleResp = types.ScheduleCourseResponse{
	Code: types.ParamInvalid,
}

var link map[string][]string

var visit map[string]int

var cy map[string]string

func dfs(x string) int {
	for _, value := range link[x] {
		if visit[value] == 0 {
			visit[value] = 1
			if cy[value] == "-1" || dfs(cy[value]) == 1 {
				cy[value] = x
				return 1
			}
		}
	}
	return 0
}

//输入老师和课程的关系，返回最大的匹配任意解（每个老师只能安排一节课，每节课只能安排一个老师）
//排课求解器，使老师绑定课程的最优解， 老师有且只能绑定一个课程
//TODO 需要操作绑定还是只返回结果即可，需要做课程绑定校验吗
func GetSchedule(c *gin.Context) {
	fmt.Println("hhh")
	var scheReq types.ScheduleCourseRequest
	if err := c.ShouldBindJSON(&scheReq); err != nil {
		fmt.Println("error1")
		c.JSON(http.StatusBadRequest, paramInvalidScheduleResp)
		return
	}

	link = scheReq.TeacherCourseRelationShip

	visit = make(map[string]int)
	cy = make(map[string]string)

	for _, v := range link {
		for _, value := range v {
			visit[value] = 0
			cy[value] = "-1"
		}
	}

	var ans int = 0

	for k, _ := range link {
		for in, _ := range visit {
			visit[in] = 0
		}
		ans += dfs(k)
	}

	cx := make(map[string]string)

	for k, v := range cy {
		if v != "-1" {
			cx[v] = k
		}
	}

	fmt.Println("error3")
	c.JSON(http.StatusOK, types.ScheduleCourseResponse{Code: types.OK, Data: cx})

}
