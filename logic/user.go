package logic

import (
	"WebBlog/dao/mysql"
	"WebBlog/models"
	"WebBlog/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//判断用户是否存在

	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	//生成UID
	userID := snowflake.GetID()

	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//构造一个User实例
	//保存进数据库
	return mysql.InsertUp(user)
}
