package member

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateMember(c *gin.Context) {
	var req types.CreateMemberRequest
	resp := &types.CreateMemberResponse{}
	defer func() {
		c.JSON(http.StatusOK, *resp)
	}()
	paramInvalidResp := types.CreateMemberResponse{
		Code: types.ParamInvalid,
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, paramInvalidResp)
	}
	// TODO 只有管理员才能添加
	// TODO 密码哈希
	if l := len(req.Nickname); l < 4 || l > 20 {
		c.JSON(http.StatusOK, paramInvalidResp)
		return
	}
	if l := len(req.Username); l < 8 || l > 20 {
		resp = &paramInvalidResp
		return
	}
	if l := len(req.Password); l < 8 || l > 20 || !validPassword(req.Password) {
		resp = &paramInvalidResp
		return
	}
	if req.UserType != types.Admin && req.UserType != types.Teacher && req.UserType != types.Student {
		resp = &paramInvalidResp
		return
	}

	member := model.Member{
		Nickname: req.Nickname,
		Username: req.Username,
		Password: req.Password,
		UserType: req.UserType,
	}
	result := database.MysqlDB.Create(&member)
	if result.Error != nil {
		panic(result.Error)
	}
	resp = &types.CreateMemberResponse{
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
