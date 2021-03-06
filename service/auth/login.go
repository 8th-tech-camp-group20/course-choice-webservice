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
		c.JSON(http.StatusOK, types.LoginResponse{Code: types.WrongPassword})
		return
	}
	//用户存在且没被删除且密码正确
	session := sessions.Default(c)
	session.Set("camp-session", user.ID)
	session.Save()
	//fmt.Println(session.Get("camp-session"))
	c.JSON(http.StatusOK, types.LoginResponse{Code: types.OK, Data: struct{ UserID string }{UserID: strconv.Itoa(int(user.ID))}})
}
