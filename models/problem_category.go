package models

import "gorm.io/gorm"

type ProblemCategory struct {
	gorm.Model
	CategoryId    uint           `gorm:"column:category_id;type:varchar(100);" json:"category_id"`
	ProblemId     uint           `gorm:"column:problem_id;type:varchar(100);" json:"problem_id"`
	CategoryBasic *CategoryBasic `gorm:"foreignKey:id;references:category_id"` //关联分类基础信息表
}

func (table *ProblemCategory) TableName() string {
	return "problem_category"
}
