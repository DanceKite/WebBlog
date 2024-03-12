package logic

import (
	"WebBlog/dao/mysql"
	"WebBlog/models"
	"WebBlog/pkg/snowflake"
	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GetID()
	// 2. 保存到数据库
	return mysql.CreatePost(p)
}

// GetPostById 根据帖子ID查询帖子详情
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合我们接口想要的数据
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}

	// 根据作者ID查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}

	// 根据社区ID查询社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return
	}

	// 接口数据拼接
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}

	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据作者ID查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		// 根据社区ID查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}

		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postDetail)
	}
	return
}
