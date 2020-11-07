package mappers

import "github.com/Masterminds/squirrel"

type UserMapper struct{}

func (*UserMapper) GetUserList() *SqlMapper {
	return Mapper(squirrel.Select("user_id", "user_name").
		From("users").
		OrderBy("user_id desc").
		Limit(10).ToSql())
}
