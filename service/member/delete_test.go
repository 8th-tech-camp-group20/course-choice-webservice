package member

import (
	"course-choice-webservice/database"
	"course-choice-webservice/types"
	"errors"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

func TestDelete(t *testing.T) {
	global := database.MysqlDB
	defer func() {
		database.MysqlDB = global
	}()
	err := global.Transaction(func(tx *gorm.DB) error {
		database.MysqlDB = tx
		// 创建
		resp := createMemberService(&types.CreateMemberRequest{
			Nickname: "66666", Username: "rust-language", Password: "Rust1Language", UserType: types.Student,
		})
		if resp.Code != types.OK {
			t.Errorf("resp.Code = %v, expect %v", resp.Code, types.OK)
		}
		userId, err := strconv.ParseUint(resp.Data.UserID, 10, 0)
		if err != nil {
			t.Errorf(err.Error())
		}
		// 删除更大 ID, 应不存在
		resp1 := deleteMemberService(&types.DeleteMemberRequest{
			UserID: strconv.FormatUint(userId+1, 10),
		})
		if resp1.Code != types.UserNotExisted {
			t.Errorf("resp1.Code = %v, expect %v", resp1.Code, types.UserNotExisted)
		}
		// 正常删除
		resp1 = deleteMemberService(&types.DeleteMemberRequest{
			UserID: resp.Data.UserID,
		})
		if resp1.Code != types.OK {
			t.Errorf("resp1.Code = %v, expect %v", resp1.Code, types.OK)
		}
		// 重复删除
		resp1 = deleteMemberService(&types.DeleteMemberRequest{
			UserID: resp.Data.UserID,
		})
		if resp1.Code != types.UserHasDeleted {
			t.Errorf("resp1.Code = %v, expect %v", resp1.Code, types.UserHasDeleted)
		}
		// 返回错误以 Rollback
		return errors.New("rollback txn")
	})
	if errStr := err.Error(); errStr != "rollback txn" {
		t.Errorf(err.Error())
	}
}
