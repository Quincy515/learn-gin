package valueobjs

type UserExtra struct {
	UserPhone string `gorm:"column:user_phone" json:"user_phone"`
	UserCity  string `gorm:"column:user_city" json:"user_city"`
	UserQq    string `gorm:"column:user_qq" json:"user_qq"`
}

func NewUserExtra(attrs ...UserExtraAttrFunc) *UserExtra {
	extra := &UserExtra{}
	UserExtraAttrFuncs(attrs).apply(extra)
	return extra
}

func (u UserExtra) Equals(other *UserExtra) bool {
	return u.UserPhone == other.UserPhone && u.UserQq == other.UserQq && u.UserCity == other.UserCity
}
