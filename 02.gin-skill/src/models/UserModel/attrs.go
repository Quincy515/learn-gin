package UserModel

type UserModelAttrFunc func(u *UserModelImpl)

type UserModelAttrFuncs []UserModelAttrFunc

func WithUserID(id int) UserModelAttrFunc {
	return func(u *UserModelImpl) {
		u.UserID = id
	}
}
func WithUserName(name string) UserModelAttrFunc {
	return func(u *UserModelImpl) {
		u.UserName = name
	}
}

func (this UserModelAttrFuncs) Apply(u *UserModelImpl) {
	for _, f := range this {
		f(u)
	}
}
