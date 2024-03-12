package logic

import (
	"WebBlog/dao/mysql"
	"WebBlog/models"
	"WebBlog/pkg/jwt"
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

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递指针，函数内部可以修改结构体的值,就能拿到user.userID
	if err = mysql.Login(user); err != nil {
		return nil, err
	}
	//生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return
}
