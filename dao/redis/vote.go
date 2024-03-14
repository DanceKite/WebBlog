package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

/* 投票的几种情况
direction=1时，有两种情况:
	1. 之前没有投过票，现在投赞成票  --> 更新分数和投票记录  差值的绝对值：1  +432
	2. 之前投反对票，现在改投赞成票  --> 更新分数和投票记录  差值的绝对值：2  +432*2
direction=0时，有两种情况:
	1. 之前投过反对票，现在要取消投票  --> 更新分数和投票记录  差值的绝对值：1  +432
	2. 之前投过赞成票，现在要取消投票  --> 更新分数和投票记录  差值的绝对值：1  -432
direction=-1时，有两种情况:
	1. 之前没有投过票，现在投反对票  --> 更新分数和投票记录  差值的绝对值：1  -432
	2. 之前投赞成票，现在改投反对票  --> 更新分数和投票记录  差值的绝对值：2  -432*2

投票的限制：
每个帖子自发表之日起一个星期之内允许投票，一个用户对一个帖子只允许投一票，超过一个星期不允许投票
	1. 到期之后将redis中保存的赞成票和反对票的数量存储到mysql表中
	2. 到期之后删除 KeyPostVotedPrefix + post_id 这个key

*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票的值多少分
)

var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
	ErrorVoted          = errors.New("已经投过票了")
)

// CreatePost 创建帖子
func CreatePost(postID, communityID int64) error {
	pipeline := rdb.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 初始化帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 把帖子id加到社区的set集合里面
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}

// VoteForPost 为帖子投票
func VoteForPost(userID, postID string, value float64) (err error) {
	// 1.判断投票的类型
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrorVoteTimeExpire
	}
	// 2和3需要放到一个事务中

	// 2.更新帖子分数
	// 先查看用户给该帖子投过票没有
	ov := rdb.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	if value == ov {
		return ErrorVoted
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //计算两次投票的差值

	pipeline := rdb.TxPipeline()

	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	if err != nil {
		return err
	}
	// 3.记录用户的投票记录
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}

	_, err = pipeline.Exec()
	return err
}
