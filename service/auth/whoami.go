package auth

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Whoami(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("camp-session")
	if uid == nil {
		c.JSON(http.StatusOK, types.WhoAmIResponse{Code: types.LoginRequired})
		return
	}
	var user model.Member
	result := database.MysqlDB.First(&user, uid)
	if result.RowsAffected != 0 {
		c.JSON(http.StatusOK, types.WhoAmIResponse{Code: types.OK, Data: struct {
			UserID   string
			Nickname string
			Username string
			UserType types.UserType
		}{UserID: strconv.Itoa(int(user.ID)), Nickname: user.Nickname, Username: user.Username, UserType: user.UserType}})
		return
	}
	//uid可能不存在或者被删除
	var count int64
	database.MysqlDB.Unscoped().Model(&model.Member{}).Where("id = ?", uid).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, types.WhoAmIResponse{Code: types.UserHasDeleted})
	} else {
		c.JSON(http.StatusOK, types.WhoAmIResponse{Code: types.UserNotExisted})
	}
	return
}
