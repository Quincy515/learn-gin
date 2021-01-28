package models

import "fmt"

// IModel 接口
type IModel interface {
	ToString() string // 方法
}

// Model 抽象类，所有都要主键和实体名称
type Model struct {
	Id   int    // 主键 - 判断实体实体之间是否相等
	Name string // 实体名称
}

// SetName 对抽象类的设定
func (m *Model) SetName(name string) {
	m.Name = name
}

func (m *Model) SetId(id int) {
	m.Id = id
}

func (m *Model) ToString() string {
	return fmt.Sprintf("Entity is: %s, id is: %d", m.Name, m.Id)
}
