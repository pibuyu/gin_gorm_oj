package models

import "gorm.io/gorm"

type CategoryBasic struct {
	gorm.Model
	Identity string `gorm:"column:identity;type:varchar(100);" json:"identity"` //分类的唯一标识
	Name     string `gorm:"column:name;type:varchar(100);" json:"name"`         //问题分类名
	ParentId int    `gorm:"column:parent_id;type:int;" json:"parent_id"`        //父级ID
}

func (table *CategoryBasic) TableName() string {
	return "category_basic"
}
