package model

import (
	"github.com/jinzhu/gorm"
)

// Channel テーブル定義
type Channel struct {
	gorm.Model
	ChannelID string `gorm:"column:channel_id"`
	Name      string `gorm:"column:name"`
	Title     string `gorm:"column:title"`
	URL       string `gorm:"column:url"`
	Thumbnail string `gorm:"column:thumbnail"`
	Category  int    `gorm:"column:category"`
	Rotation  *int   `gorm:"column:rotation"`
}
