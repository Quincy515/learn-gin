package valueobjs

type UserExtraAttrFunc func(model *UserExtra) // 设置 User 属性的 函数类型
type UserExtraAttrFuncs []UserExtraAttrFunc

// 传参数
func WithUserPhone(phone string) UserExtraAttrFunc {
	return func(u *UserExtra) {
		u.UserPhone = phone
	}
}

func WithUserQQ(qq string) UserExtraAttrFunc {
	return func(u *UserExtra) {
		u.UserQq = qq
	}
}

func WithUserCity(city string) UserExtraAttrFunc {
	return func(u *UserExtra) {
		u.UserCity = city
	}
}

// apply 方法 循环 UserExtraAttrFuncs 内容执行函数
func (u UserExtraAttrFuncs) apply(model *UserExtra) {
	for _, f := range u {
		f(model)
	}
}
