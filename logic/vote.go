package logic

import (
	"WebBlog/dao/redis"
	"WebBlog/models"
	"go.uber.org/zap"
	"strconv"
)

// 基于用户投票的相关算法 https://ruanyifeng.com/blog/algorithm/

// 在本项目使用简化版的投票分数

/* 投票的几种情况
direction=1时，有两种情况:
	1. 之前没有投过票，现在投赞成票
	2. 之前投反对票，现在改投赞成票
direction=0时，有两种情况:
	1. 之前投过赞成票，现在要取消投票
	2. 之前投过反对票，现在要取消投票
direction=-1时，有两种情况:
	1. 之前没有投过票，现在投反对票
	2. 之前投赞成票，现在改投反对票

投票的限制：
每个帖子自发表之日起一个星期之内允许投票，一个用户对一个帖子只允许投一票，超过一个星期不允许投票
	1. 到期之后将redis中保存的赞成票和反对票的数量存储到mysql表中
	2. 到期之后删除 KeyPostVotedPrefix + post_id 这个key

*/

// VoteForPost 为帖子投票
func VoteForPost(userID int64, p *models.ParamVoteData) (err error) {
	zap.L().Debug("VoteForPost", zap.Any("userID", userID), zap.Any("postID", p.PostID), zap.Any("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
