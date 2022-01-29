package member

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var paramInvalidUpdateMemberResp = types.UpdateMemberResponse{
	Code: types.ParamInvalid,
}

func UpdateMember(c *gin.Context) {
	var req types.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, paramInvalidUpdateMemberResp)
		return
	}
	resp := updateMemberService(&req)
	c.JSON(http.StatusOK, *resp)
}

func updateMemberService(req *types.UpdateMemberRequest) *types.UpdateMemberResponse {
	// 检验昵称
	if l := len(req.Nickname); l < 4 || l > 20 {
		return &paramInvalidUpdateMemberResp
	}
	// 构建查询条件
	userId, err := strconv.ParseUint(req.UserID, 10, 0)
	if err != nil {
		return &paramInvalidUpdateMemberResp
	}
	userToUpd := model.Member{
		Model: gorm.Model{
			ID: uint(userId),
		},
	}
	database.MysqlDB.Model(&userToUpd).Update("nickname", req.Nickname)
	return &types.UpdateMemberResponse{
		Code: types.OK,
	}
}
