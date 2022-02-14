package member

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	paramInvalidGetMemberResp     = types.GetMemberResponse{Code: types.ParamInvalid}
	paramInvalidGetMemberListResp = types.GetMemberListResponse{Code: types.ParamInvalid}
)

func GetMember(c *gin.Context) {
	var req types.GetMemberRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, paramInvalidGetMemberResp)
		return
	}
	resp := getMemberService(&req)
	c.JSON(http.StatusOK, *resp)
}

func getMemberService(req *types.GetMemberRequest) *types.GetMemberResponse {
	userId, err := strconv.ParseUint(req.UserID, 10, 0)
	if err != nil {
		return &paramInvalidGetMemberResp
	}
	apiUser := struct {
		ID                 uint
		Nickname, Username string
		UserType           types.UserType
	}{}
	result := database.MysqlDB.Model(&model.Member{}).First(&apiUser, uint(userId))
	if result.RowsAffected == 0 {
		var count int64
		database.MysqlDB.Unscoped().Model(&model.Member{}).Where("id = ?", userId).Count(&count)
		if count > 0 {
			return &types.GetMemberResponse{Code: types.UserHasDeleted}
		} else {
			return &types.GetMemberResponse{Code: types.UserNotExisted}
		}
	}
	return &types.GetMemberResponse{
		Code: types.OK,
		Data: types.TMember{
			UserID:   strconv.FormatUint(uint64(apiUser.ID), 10),
			Nickname: apiUser.Nickname,
			Username: apiUser.Username,
			UserType: apiUser.UserType,
		},
	}
}

func GetMemberList(c *gin.Context) {
	var req types.GetMemberListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, paramInvalidGetMemberListResp)
		return
	}
	resp := getMemberListService(&req)
	c.JSON(http.StatusOK, *resp)
}

func getMemberListService(req *types.GetMemberListRequest) *types.GetMemberListResponse {
	var apiUsers []struct {
		ID                 uint
		Nickname, Username string
		UserType           types.UserType
	}
	database.MysqlDB.Model(&model.Member{}).Order("id").Limit(req.Limit).Offset(req.Offset).Find(&apiUsers)
	resp := &types.GetMemberListResponse{
		Code: types.OK,
	}
	for i := range apiUsers {
		userID := strconv.FormatUint(uint64(apiUsers[i].ID), 10)
		resp.Data.MemberList = append(resp.Data.MemberList, types.TMember{
			UserID:   userID,
			Nickname: apiUsers[i].Nickname,
			Username: apiUsers[i].Username,
			UserType: apiUsers[i].UserType,
		})
	}
	return resp
}
