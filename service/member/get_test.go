package member

import (
	"course-choice-webservice/database"
	"course-choice-webservice/types"
	"errors"
	"gorm.io/gorm"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		nickname, username, password string
		usertype                     types.UserType
	}{
		{"66666", "rust-language", "Rust1Language", types.Student},
		{"88888", "go-language", "Go1Language", types.Student},
		{"99999", "java-language", "Java1Language", types.Student},
	}
	global := database.MysqlDB
	defer func() {
		database.MysqlDB = global
	}()
	err := global.Transaction(func(tx *gorm.DB) error {
		database.MysqlDB = tx
		userIds := make([]string, 0, 3)
		for i := range tests {
			resp := createMemberService(&types.CreateMemberRequest{
				Nickname: tests[i].nickname,
				Username: tests[i].username,
				Password: tests[i].password,
				UserType: tests[i].usertype,
			})
			if resp.Code != types.OK {
				t.Errorf("resp.Code = %v, expect %v", resp.Code, types.OK)
			}
			userIds = append(userIds, resp.Data.UserID)
		}
		for i := range tests {
			resp := getMemberService(&types.GetMemberRequest{UserID: userIds[i]})
			if resp.Code != types.OK {
				t.Errorf("resp.Code = %v, expect %v", resp.Code, types.OK)
			}
			if resp.Data.UserID != userIds[i] {
				t.Errorf("resp.Data.UserID = %v, expect %v", resp.Data.UserID, userIds[i])
			}
			if resp.Data.Username != tests[i].username {
				t.Errorf("resp.Data.Username = %v, expect %v", resp.Data.Username, tests[i].username)
			}
			if resp.Data.Nickname != tests[i].nickname {
				t.Errorf("resp.Data.Nickname = %v, expect %v", resp.Data.Nickname, tests[i].nickname)
			}
			if resp.Data.UserType != tests[i].usertype {
				t.Errorf("resp.Data.UserType = %v, expect %v", resp.Data.UserType, tests[i].usertype)
			}
		}
		// 返回错误以 Rollback
		return errors.New("rollback txn")
	})
	if errStr := err.Error(); errStr != "rollback txn" {
		t.Errorf(err.Error())
	}
}

func TestGetList(t *testing.T) {
	tests := []struct {
		nickname, username, password string
		usertype                     types.UserType
	}{
		{"66666", "rust-language", "Rust1Language", types.Student},
		{"88888", "go-language", "Go1Language", types.Student},
		{"99999", "java-language", "Java1Language", types.Student},
	}
	global := database.MysqlDB
	defer func() {
		database.MysqlDB = global
	}()
	err := global.Transaction(func(tx *gorm.DB) error {
		database.MysqlDB = tx
		userIds := make([]string, 0, 3)
		for i := range tests {
			resp := createMemberService(&types.CreateMemberRequest{
				Nickname: tests[i].nickname,
				Username: tests[i].username,
				Password: tests[i].password,
				UserType: tests[i].usertype,
			})
			if resp.Code != types.OK {
				t.Errorf("resp.Code = %v, expect %v", resp.Code, types.OK)
			}
			userIds = append(userIds, resp.Data.UserID)
		}
		compareResult := func(offset, limit, start, end int) {
			resp := getMemberListService(&types.GetMemberListRequest{
				Offset: offset,
				Limit:  limit,
			})
			list := resp.Data.MemberList
			idx := 0
			for i := start; i < end; i++ {
				if list[idx].Username != tests[i].username {
					t.Errorf("list[i].Username = %v, expect %v", list[i].Username, tests[i].username)
				}
				if list[idx].Nickname != tests[i].nickname {
					t.Errorf("list[i].Nickname = %v, expect %v", list[i].Nickname, tests[i].nickname)
				}
				if list[idx].UserType != tests[i].usertype {
					t.Errorf("list[i].UserType = %v, expect %v", list[i].UserType, tests[i].usertype)
				}
				if list[idx].UserID != userIds[i] {
					t.Errorf("list[i].UserID = %v, expect %v", list[i].UserID, userIds[i])
				}
				idx++
			}
		}
		compareResult(0, 2, 0, 2)
		compareResult(2, 2, 2, 3)
		compareResult(1, 1, 1, 2)
		// 返回错误以 Rollback
		return errors.New("rollback txn")
	})
	if errStr := err.Error(); errStr != "rollback txn" {
		t.Errorf(err.Error())
	}
}
