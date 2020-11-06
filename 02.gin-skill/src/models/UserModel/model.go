package UserModel

type UserModelImpl struct {
	UserID   int
	UserName string
}

func New() *UserModelImpl {
	return &UserModelImpl{}
}
