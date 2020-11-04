package src

type Topic struct {
	TopicID         int    `json:"id"`
	TopicTitle      string `json:"title"  binding:"min=4,max=20" `
	TopicShortTitle string `json:"stitle"  binding:"nefield=TopicTitle"`
	UserIP          string `json:"ip" binding:"ip4_addr"`
	TopicScore      int    `json:"score" binding:"omitempty,gt=5"`
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
