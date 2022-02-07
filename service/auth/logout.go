package auth

import (
	"course-choice-webservice/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	//fmt.Println(session.Get("camp-session"))
	session.Delete("camp-session")
	session.Save()
	//fmt.Println(session.Get("camp-session"))
	c.JSON(http.StatusOK, types.LogoutResponse{Code: types.OK})
}
