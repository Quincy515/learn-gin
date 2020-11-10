package models

type Source struct {
	SourceID     int `gorm:"column:source_id"`
	SourceName   string
}
