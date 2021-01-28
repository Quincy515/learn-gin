package models

type UserLogAttrFunc func(model *UserLogModel)
type UserLogAttrFuncs []UserLogAttrFunc

func WithUserLogType(t uint8) UserLogAttrFunc {
	return func(u *UserLogModel) {
		u.LogType = t
	}
}

func WithUserLogComment(comment string) UserLogAttrFunc {
	return func(u *UserLogModel) {
		u.LogComment = comment
	}
}

func (u UserLogAttrFuncs) apply(model *UserLogModel) {
	for _, f := range u {
		f(model)
	}
}
