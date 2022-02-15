package service

import (
	"course-choice-webservice/service/auth"
	"course-choice-webservice/service/course"
	"course-choice-webservice/service/member"
	"course-choice-webservice/service/student"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// 成员管理
	g.POST("/member/create", member.CreateMember)
	g.GET("/member", member.GetMember)
	g.GET("/member/list", member.GetMemberList)
	g.POST("/member/update", member.UpdateMember)
	g.POST("/member/delete", member.DeleteMember)

	// 登录

	g.POST("/auth/login", auth.Login)
	g.POST("/auth/logout", auth.Logout)
	g.GET("/auth/whoami", auth.Whoami)

	// 排课
	g.POST("/course/create", course.CreateCourse)
	g.GET("/course/get", course.GetCourse)

	g.POST("/teacher/bind_course", course.BindCourse)
	g.POST("/teacher/unbind_course", course.UnbindCourse)
	g.GET("/teacher/get_course", course.GetTeacherCourse)
	g.POST("/course/schedule", course.GetSchedule)

	// 抢课
	g.POST("/student/book_course", student.BookCourse)
	g.GET("/student/course", student.GetStudentCourse)

}
