package auth

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"crypto/sha256"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	var req = types.LoginRequest{}
	//解析json到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, types.LoginResponse{Code: types.ParamInvalid})
		return
	}
	//验证密码
	var password [32]byte = sha256.Sum256([]byte(req.Password + req.Password))
	var user model.Member
	result := database.MysqlDB.Where(&model.Member{Username: req.Username, Password: password[:]}).First(&user)
	if result.Error != nil {
		//TODO 这里不清楚具体需求是要区分开用户不存在，用户已删除，密码错误还是统一返回密码错误
		c.JSON(http.StatusOK, types.LoginResponse{Code: types.WrongPassword})
		return
	}
	//用户存在且没被删除且密码正确
	c.JSON(http.StatusOK, types.LoginResponse{Code: types.OK, Data: struct{ UserID string }{UserID: string(user.ID)}})
	
}
