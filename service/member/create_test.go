package member

import (
	"course-choice-webservice/database"
	"course-choice-webservice/types"
	"errors"
	"gorm.io/gorm"
	"testing"
)

func TestCreate(t *testing.T) {
	global := database.MysqlDB
	defer func() {
		database.MysqlDB = global
	}()
	err := global.Transaction(func(tx *gorm.DB) error {
		database.MysqlDB = tx
		// 昵称, 用户名, 密码, 用户类型校验
		resp := createMemberService(&types.CreateMemberRequest{})
		if resp.Code != types.ParamInvalid {
			t.Errorf("resp.Code = %v, expect %v", resp.Code, types.ParamInvalid)
		}
		resp = createMemberService(&types.CreateMemberRequest{Nickname: "66666"})
		if resp.Code != types.ParamInvalid {
			t.Errorf("resp.Code = %v, expect %v", resp.Code, types.ParamInvalid)
		}
		resp = createMemberService(&types.CreateMemberRequest{Nickname: "66666", Username: "rust-language"})
		if resp.Code != types.ParamInvalid {
			t.Errorf("resp.Code = %v, expect %v", resp.Code, types.ParamInvalid)
		}
		resp = createMemberService(&types.CreateMemberRequest{
			Nickname: "66666", Username: "rust-language", Password: "rust-language",
		})
		if resp.Code != types.ParamInvalid {
			t.Errorf("resp.Code = %v, expect %v", resp.Code, types.ParamInvalid)
		}
		resp = createMemberService(&types.CreateMemberRequest{
			Nickname: "66666", Username: "rust-language", Password: "Rust1Language",
		})
		if resp.Code != types.ParamInvalid {
			t.Errorf("resp.Code = %v, expect %v", resp.Code, types.ParamInvalid)
		}
		resp = createMemberService(&types.CreateMemberRequest{
			Nickname: "66666", Username: "rust-language", Password: "Rust1Language", UserType: types.Student,
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
