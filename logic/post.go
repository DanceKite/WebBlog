package logic

import (
	"WebBlog/dao/mysql"
	"WebBlog/dao/redis"
	"WebBlog/models"
	"WebBlog/pkg/snowflake"
	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GetID()
	// 2. 保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		zap.L().Error("mysql.CreatePost() failed", zap.Error(err))
		return
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
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

// GetPostList2 获取帖子列表升级版
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 根据前端传来的参数动态获得帖子列表（按分数、按时间）
	// 1.获取分页参数
	// 2.去redis查询获取id列表
	// 3.根据id去mysql查询帖子的详细信息
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		zap.L().Error("redis.GetPostIDsInOrder failed", zap.Error(err))
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder return 0 data")
		return
	}
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostListByIDs failed", zap.Error(err))
		return
	}
	// 根据ids提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
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
			VoteNum:         voteData[idx],
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 根据前端传来的参数动态获得帖子列表（按分数、按时间）
	// 1.获取分页参数
	// 2.去redis查询获取id列表
	// 3.根据id去mysql查询帖子的详细信息
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		zap.L().Error("redis.GetPostIDsInOrder failed", zap.Error(err))
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder return 0 data")
		return
	}
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostListByIDs failed", zap.Error(err))
		return
	}
	// 根据ids提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
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
			VoteNum:         voteData[idx],
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postDetail)
	}
	return
}

// GetPostListNew 获取帖子列表(将两个查询逻辑合并)
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		// 查询所有的帖子
		data, err = GetPostList2(p)
	} else {
		// 查询社区的帖子
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return
	}
	return
}
