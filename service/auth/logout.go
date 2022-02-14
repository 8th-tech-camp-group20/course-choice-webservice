package auth

import (
	"course-choice-webservice/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("camp-session")
	if uid == nil {
		//未登录
		c.JSON(http.StatusOK, types.LogoutResponse{Code: types.LoginRequired})
		return
	}
	//fmt.Println(session.Get("camp-session"))
	session.Delete("camp-session")
	session.Save()
	//fmt.Println(session.Get("camp-session"))
	c.JSON(http.StatusOK, types.LogoutResponse{Code: types.OK})
}
