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

var paramInvalidDeleteMemberResp = types.DeleteMemberResponse{
	Code: types.ParamInvalid,
}

func DeleteMember(c *gin.Context) {
	var req types.DeleteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, paramInvalidDeleteMemberResp)
		return
	}
	resp := deleteMemberService(&req)
	c.JSON(http.StatusOK, *resp)
}

func deleteMemberService(req *types.DeleteMemberRequest) *types.DeleteMemberResponse {
	userId, err := strconv.ParseUint(req.UserID, 10, 0)
	if err != nil {
		return &paramInvalidDeleteMemberResp
	}
	rows := database.MysqlDB.Delete(&model.Member{Model: gorm.Model{
		ID: uint(userId),
	}}).RowsAffected
	if rows == 0 {
		var count int64
		database.MysqlDB.Unscoped().Model(&model.Member{}).Where("id = ?", userId).Count(&count)
		if count > 0 {
			return &types.DeleteMemberResponse{Code: types.UserHasDeleted}
		} else {
			return &types.DeleteMemberResponse{Code: types.UserNotExisted}
		}
	}
	return &types.DeleteMemberResponse{Code: types.OK}
}
