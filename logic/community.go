package logic

import (
	"WebBlog/dao/mysql"
	"WebBlog/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查询到所有的社区（community_id, community_name）以列表的形式返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
