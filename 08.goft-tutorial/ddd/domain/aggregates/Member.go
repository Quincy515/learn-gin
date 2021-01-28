package aggregates

import "goft-tutorial/ddd/domain/models"

// Member 会员聚合 -- 会员：用户+日志+...组成
type Member struct {
	User *models.UserModel
	Log  *models.UserLogModel
	// 充值、社交、隐私信息
}
