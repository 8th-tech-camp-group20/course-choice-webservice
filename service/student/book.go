package student

import (
	"course-choice-webservice/database"
	"course-choice-webservice/model"
	"course-choice-webservice/types"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var (
	paramInvalidBookCourseRequest = types.BookCourseResponse{Code: types.ParamInvalid}
)

func BookCourse(c *gin.Context) {
	var req types.BookCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := bookCourseService(&req)
	c.JSON(http.StatusOK, *resp)
}

func bookCourseService(req *types.BookCourseRequest) *types.BookCourseResponse {
	studentId, err := strconv.ParseUint(req.StudentID, 10, 64)
	if err != nil {
		return &paramInvalidBookCourseRequest
	}
	courseId := req.CourseID

	rc := database.RedisClient.Get()
	defer rc.Close()

	//flag, err := redis.Int(rc.Do("GET", courseId+"flag"))
	//if err != nil {
	//	fmt.Println("redis get failed:", err)
	//}
	//if flag == 0 {
	//	return &types.BookCourseResponse{
	//		Code: types.CourseNotAvailable,
	//	}
	//}

	ex, err := redis.Bool(rc.Do("SISMEMBER", "student"+req.StudentID, courseId))
	if err != nil {
		fmt.Println("redis sadd failed:", err)
	}

	if ex {
		return &types.BookCourseResponse{
			Code: types.StudentHasCourse,
		}
	}

	//判断redis中是否有key 为 courseId 的项
	ex, err = redis.Bool(rc.Do("EXISTS", courseId))
	if err != nil {
		fmt.Println("redis get exist failed:", err)
	}

	if !ex {
		//redis中没有key 为 courseId 的项就从数据库读出其对应的容量
		thisCourse := model.Course{}
		rows := database.MysqlDB.Table("course").Where("id=?", courseId).Find(&thisCourse).RowsAffected
		//查询的信息为空，返回课程不存在
		if rows == 0 {
			return &types.BookCourseResponse{
				Code: types.CourseNotExisted,
			}
		} else {
			//容量写入redis
			_, err = rc.Do("SETNX", courseId, thisCourse.Cap)
			if err != nil {
				fmt.Println("redis set failed:", err)
			}
		}
	}
	//redis中有key 为 courseId 的项就读出courseId对应的容量剩余数
	count, err := redis.Int(rc.Do("GET", courseId))
	if err != nil {
		fmt.Println("redis get failed:", err)
	}
	count, err = redis.Int(rc.Do("DECR", courseId))
	if err != nil {
		fmt.Println("redis decr failed:", err)
	}
	//剩余数<0，返回抢课失败
	if count < 0 {
		_, err = rc.Do("INCR", courseId)
		if err != nil {
			fmt.Println("redis incr failed:", err)
		}
		_, err = rc.Do("SET", courseId+"flag", 0)
		if err != nil {
			fmt.Println("redis set failed:", err)
		}
		return &types.BookCourseResponse{
			Code: types.CourseNotAvailable,
		}
	} else {
		//根据studentId查信息
		stu := model.Member{
			Model: gorm.Model{
				ID: uint(studentId),
			},
		}
		rows := database.MysqlDB.First(&stu, "id = ?", studentId).RowsAffected
		//查询的信息为空，返回学生不存在
		if rows == 0 {
			_, err = rc.Do("INCR", courseId)
			if err != nil {
				fmt.Println("redis incr failed:", err)
			}
			return &types.BookCourseResponse{
				Code: types.StudentNotExisted,
			}
		}
		////该studentId对应type不为student，返回无权限
		if stu.UserType != types.Student {
			_, err = rc.Do("INCR", courseId)
			if err != nil {
				fmt.Println("redis incr failed:", err)
			}
			return &types.BookCourseResponse{
				Code: types.PermDenied,
			}
		}
		bookSuccess(studentId, courseId)
		_, err = rc.Do("SADD", "student"+req.StudentID, courseId)
		if err != nil {
			fmt.Println("redis sadd failed:", err)
		}
		return &types.BookCourseResponse{
			Code: types.OK,
		}
	}
}

func bookSuccess(studentId uint64, courseId string) {
	//查询这个学生id是否存在
	stuCou := model.StudentCourse{
		StudentId: studentId,
	}
	rows := database.MysqlDB.Table("book_course").Where("student_id = ?", studentId).Find(&stuCou).RowsAffected

	//存在
	if rows != 0 {
		courseIdList := stuCou.CourseIdList
		newCourseIdList := courseIdList + "," + courseId
		stuCou.CourseIdList = newCourseIdList
		database.MysqlDB.Table("book_course").Where("student_id = ?", studentId).Updates(stuCou)
	} else {
		//不存在
		stuCou = model.StudentCourse{
			StudentId:    studentId,
			CourseIdList: courseId,
		}
		database.MysqlDB.Table("book_course").Create(&stuCou)
	}
}
