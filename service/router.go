package service

import (
	"course-choice-webservice/service/course"
	"course-choice-webservice/service/member"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1") //TODO 建议这里改成空的，不然直接访问下面的接口访问不到

	// 成员管理
	g.POST("/member/create", member.CreateMember)
	g.GET("/member", member.GetMember)
	g.GET("/member/list", member.GetMemberList)
	g.POST("/member/update", member.UpdateMember)
	g.POST("/member/delete", member.DeleteMember)

	// 登录

	g.POST("/auth/login")
	g.POST("/auth/logout")
	g.GET("/auth/whoami")

	// 排课
	g.POST("/course/create", course.CreateCourse)
	g.GET("/course/get", course.GetCourse)

	g.POST("/teacher/bind_course", course.BindCourse)
	g.POST("/teacher/unbind_course", course.UnbindCourse)
	g.GET("/teacher/get_course", course.GetTeacherCourse)
	g.POST("/course/schedule")

	// 抢课
	g.POST("/student/book_course")
	g.GET("/student/course")

}
