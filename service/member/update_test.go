package member

import (
	"course-choice-webservice/database"
	"course-choice-webservice/types"
	"errors"
	"gorm.io/gorm"
	"testing"
)

func TestUpdate(t *testing.T) {
	global := database.MysqlDB
	defer func() {
		database.MysqlDB = global
	}()
	err := global.Transaction(func(tx *gorm.DB) error {
		database.MysqlDB = tx
		// 昵称校验
		resp := updateMemberService(&types.UpdateMemberRequest{})
		if resp.Code != types.ParamInvalid {
			t.Errorf("resp.Code = %v, expect %v", resp.Code, types.ParamInvalid)
		}
		// 创建并修改
		resp1 := createMemberService(&types.CreateMemberRequest{
			Nickname: "66666", Username: "rust-language", Password: "Rust1Language", UserType: types.Student,
		})
		if resp1.Code != types.OK {
			t.Errorf("resp1.Code = %v, expect %v", resp1.Code, types.OK)
		}
		resp = updateMemberService(&types.UpdateMemberRequest{
			UserID:   resp1.Data.UserID,
			Nickname: "999999",
		})
		if resp.Code != types.OK {
			t.Errorf("resp.Code = %v, expect %v", resp.Code, types.OK)
		}
		// 返回错误以 Rollback
		return errors.New("rollback txn")
	})
	if errStr := err.Error(); errStr != "rollback txn" {
		t.Errorf(err.Error())
	}
}
