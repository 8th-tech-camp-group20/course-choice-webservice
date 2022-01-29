package member

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"crypto/sha256"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var paramInvalidCreateMemberResp = types.CreateMemberResponse{
	Code: types.ParamInvalid,
}

func CreateMember(c *gin.Context) {
	var req types.CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, paramInvalidCreateMemberResp)
		return
	}
	resp := createMemberService(&req)
	c.JSON(http.StatusOK, *resp)
}

func createMemberService(req *types.CreateMemberRequest) *types.CreateMemberResponse {
	// TODO 只有管理员才能添加
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
		panic(result.Error)
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
