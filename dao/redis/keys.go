package redis

// redis key

// redis key 注意使用命名空间的方式，方便查询和拆分
const (
	keyPrefix          = "bluebell:"
	KeyPostTimeZSet    = "post:time"   //ZSet;帖子及发帖时间
	KeyPostScoreZSet   = "post:score"  //ZSet;帖子及投票分数
	KeyPostVotedZSetPF = "post:voted:" //ZSet;记录用户及投票类型;参数是post id

	KeyCommunitySetPF = "community:" //set;保存每个社区下帖子id
)

// 给Redis的key加上前缀
func getRedisKey(key string) string {
	return keyPrefix + key
}
