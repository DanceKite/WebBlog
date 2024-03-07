package mysql

import (
	"WebBlog/models"
	"crypto/md5"
	"encoding/hex"
	"errors"
)

const secret = "github.com/DanceKite"

// CheckUserExist 判断用户是否存在
func CheckUserExist(username string) (err error) {
	//执行SQL语句查询用户信息是否存在
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已经存在")
	}

	return
}

// InsertUp 向数据库中插入一条新的用户记录
func InsertUp(user *models.User) (err error) {
	//对密码进行加密
	user.Password = encryptPassword(user.Password)

	//执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func encryptPassword(oPassword string) string {
	//对密码进行加密
	h := md5.New()
	h.Write([]byte(secret))

	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
