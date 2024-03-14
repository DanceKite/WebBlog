package redis

import (
	"WebBlog/models"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1

	return rdb.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询好每篇帖子的投票数
func GetPostVoteData(ids []string) (data []int64, err error) {
	/*data = make([]int64, 0, len(ids))
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		// 查找key中值为1的数量（赞成票）
		v := rdb.ZCount(key, "1", "1").Val()
		data = append(data, v)
	}*/

	// 使用pipeline一次发送多个命令，减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}

	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}

	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()

		data = append(data, v)
	}

	return
}

// GetCommunityPostIDsInOrder 根据社区ID获取帖子id列表
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 使用zinterstore做交集
	// 针对新的zset 按之前的逻辑取数据
	// 社区的key
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))

	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}

	return getIDsFromKey(key, p.Page, p.Size)
}
