package mysql

import "WebBlog/models"

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	// sql语句
	sqlStr := `insert into post(post_id, author_id, community_id, title, content) values(?,?,?,?,?)`
	// 执行
	_, err = db.Exec(sqlStr, p.ID, p.AuthorID, p.CommunityID, p.Title, p.Content)
	return
}

// GetPostById 根据帖子ID查询帖子详情
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	// sql语句
	sqlStr := `select post_id, author_id, community_id, title, content, create_time from post where post_id = ?`
	// 执行
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, author_id, community_id, title, content, create_time from post limit ?,?`

	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}
