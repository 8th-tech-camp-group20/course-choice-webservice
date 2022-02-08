package auth

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"crypto/sha256"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Login(c *gin.Context) {
	var req = types.LoginRequest{}
	//解析json到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, types.LoginResponse{Code: types.ParamInvalid})
		return
	}
	//验证密码
	var password [32]byte = sha256.Sum256([]byte(req.Password + req.Username))
	var user model.Member
	result := database.MysqlDB.Where(&model.Member{Username: req.Username, Password: password[:]}).First(&user)
	if result.Error != nil {
		//TODO 这里不清楚具体需求是要区分开用户不存在，用户已删除，密码错误还是统一返回密码错误
		result = database.MysqlDB.Where(&model.Member{Username: req.Username}).First(&user)
		if result.RowsAffected != 0 {
			c.JSON(http.StatusOK, types.LoginResponse{Code: types.WrongPassword})
		}
		var count int64
		database.MysqlDB.Unscoped().Model(&model.Member{}).Where(&model.Member{Username: req.Username}).Count(&count)
		if count > 0 {
			c.JSON(http.StatusOK, types.WhoAmIResponse{Code: types.UserHasDeleted})
		} else {
			c.JSON(http.StatusOK, types.WhoAmIResponse{Code: types.UserNotExisted})
		}
		return
	}
	//用户存在且没被删除且密码正确
	session := sessions.Default(c)
	session.Set("camp-session", user.ID)
	session.Save()
	//fmt.Println(session.Get("camp-session"))
	c.JSON(http.StatusOK, types.LoginResponse{Code: types.OK, Data: struct{ UserID string }{UserID: strconv.Itoa(int(user.ID))}})
}
