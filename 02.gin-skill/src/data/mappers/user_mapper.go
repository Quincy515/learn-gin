package mappers

import "github.com/Masterminds/squirrel"

type UserMapper struct{}

func (*UserMapper) GetUserList() *SqlMapper {
	return Mapper(squirrel.Select("user_id", "user_name").
		From("users").
		OrderBy("user_id desc").
		Limit(10).ToSql())
}

func (*UserMapper) GetUserDetail(id int) *SqlMapper {
	return Mapper(squirrel.Select("user_id", "user_name").
		From("users").Where("user_id=?", id).ToSql())
}
