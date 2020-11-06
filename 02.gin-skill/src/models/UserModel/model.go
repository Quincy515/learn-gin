package UserModel

type UserModelImpl struct {
	UserID      int    `json:"id" form:"id"`
	UserName    string `json:"name" form:"name" binding:"required,UserName"`
	UserPwd     string `json:"user_pwd" binding:"required,min=4"`
	UserAddtime string `json:"addtime"`
}

func (*UserModelImpl) TableName() string {
	return "users"
}

// New 初始化实例
func New(attrs ...UserModelAttrFunc) *UserModelImpl {
	u := &UserModelImpl{}
	// 对 u 里每个属性进行初始化
	// 强制类型转化。
	UserModelAttrFuncs(attrs).Apply(u)
	return u
}

// Mutate 修改实例属性
func (this *UserModelImpl) Mutate(attrs ...UserModelAttrFunc) *UserModelImpl {
	UserModelAttrFuncs(attrs).Apply(this)
	return this
}
