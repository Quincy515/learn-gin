package src

// Topic 单个帖子实体
type Topic struct {
	TopicID         int    `json:"id"`
	TopicTitle      string `json:"title" binding:"min=4,max=20"`
	TopicShortTitle string `json:"stitle" binding:"required,nefield=TopicTitle"`
	TopicUrl        string `json:"url" binding:"omitempty,topicurl"`
	UserIP          string `json:"ip" binding:"ipv4"`
	TopicScore      int    `json:"score" binding:"omitempty,gt=5"`
}

// Topics 多条帖子实体
type Topics struct {
	TopicList     []Topic `json:"topics" binding:"gt=0,lt=3,topics,dive"`
	TopicListSize int     `json:"size"`
}

// CreateTopic 临时创建实体
func CreateTopic(id int, title string) Topic {
	//return Topic{id, title}
	return Topic{}
}

type TopicQuery struct {
	UserName string `json:"username" form:"username"`
	Page     int    `json:"page" form:"page" binding:"required"`
	PageSize int    `json:"pagesize" form:"pagesize"`
}
