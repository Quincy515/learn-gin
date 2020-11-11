package models

import "fmt"

type UserModel struct {
	UserID   int
	UserName string
}

func (this *UserModel) String() string {
	return fmt.Sprintf("user_id: %d, user_name: %s", this.UserID, this.UserName)
}
