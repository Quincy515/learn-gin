package mappers

import (
	"ginskill/src/models/LogModel"
	"github.com/Masterminds/squirrel"
)

type LogMapper struct{}

func (*LogMapper) AddLog(log *LogModel.LogImpl) *SqlMapper {
	return Mapper(squirrel.Insert(log.TableName()).
		Columns("log_name", "log_date").
		Values(log.LogName, log.LogDate).ToSql())
}
