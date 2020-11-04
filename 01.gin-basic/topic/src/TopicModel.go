package src

type Topic struct {
	TopicID    int    `json:"id"`
	TopicTitle string `json:"title"`
}

// CreateTopic 临时创建实体
func CreateTopic(id int, title string) Topic {
	return Topic{id, title}
}

type TopicQuery struct {
	UserName string `json:"username" form:"username"`
	Page     int    `json:"page" form:"page" binding:"required"`
	PageSize int    `json:"pagesize" form:"pagesize"`
}
