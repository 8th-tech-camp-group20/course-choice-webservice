package member

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"crypto/sha256"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

var (
	paramInvalidCreateMemberResp = types.CreateMemberResponse{
		Code: types.ParamInvalid,
	}
	unknownCreateMemberResp = types.CreateMemberResponse{Code: types.UnknownError}
)

func CreateMember(c *gin.Context) {
	var req types.CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, paramInvalidCreateMemberResp)
		return
	}
	// 获取用户 ID 以判断权限
	session := sessions.Default(c)
	uid := session.Get("camp-session")
	if uid == nil {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.LoginRequired})
		return
	}
	if !isAdmin(uid) {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.PermDenied})
		return
	}
	resp := createMemberService(&req)
	c.JSON(http.StatusOK, *resp)
}

func createMemberService(req *types.CreateMemberRequest) *types.CreateMemberResponse {
	if l := len(req.Nickname); l < 4 || l > 20 {
		return &paramInvalidCreateMemberResp
	}
	if l := len(req.Username); l < 8 || l > 20 {
		return &paramInvalidCreateMemberResp
	}
	if l := len(req.Password); l < 8 || l > 20 || !validPassword(req.Password) {
		return &paramInvalidCreateMemberResp
	}
	if req.UserType != types.Admin && req.UserType != types.Teacher && req.UserType != types.Student {
		return &paramInvalidCreateMemberResp
	}

	// 用户名不能修改，因此以用户名作为盐值
	hashedPassword := sha256.Sum256([]byte(req.Password + req.Username))
	member := model.Member{
		Nickname: req.Nickname,
		Username: req.Username,
		Password: hashedPassword[:],
		UserType: req.UserType,
	}
	result := database.MysqlDB.Create(&member)
	if result.Error != nil {
		if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				// 用户已存在
				return &types.CreateMemberResponse{Code: types.UserHasExisted}
			} else {
				return &unknownCreateMemberResp
			}
		} else {
			return &unknownCreateMemberResp
		}
	}
	return &types.CreateMemberResponse{
		Code: types.OK,
		Data: struct{ UserID string }{UserID: strconv.FormatUint(uint64(member.ID), 10)},
	}
}

func validPassword(p string) bool {
	hasUpper, hasLower, hasNumber := false, false, false
	for i := range p {
		switch {
		case p[i] >= 'a' && p[i] <= 'z':
			hasLower = true
		case p[i] >= 'A' && p[i] <= 'Z':
			hasUpper = true
		case p[i] >= '0' && p[i] <= '9':
			hasNumber = true
		}
	}
	return hasUpper && hasLower && hasNumber
}

// isAdmin 根据 Cookie 信息判断用户身份
func isAdmin(cookie interface{}) bool {
	apiUser := struct {
		UserType types.UserType
	}{}
	result := database.MysqlDB.Model(&model.Member{}).Find(&apiUser, cookie)
	if result.RowsAffected == 0 || apiUser.UserType != types.Admin {
		return false
	}
	return true
}
