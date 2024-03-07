package models

import "gorm.io/gorm"

type ProblemCategory struct {
	gorm.Model
	//CategoryId和ProblemId我自己用的gorm的类型是varchar，不对的话可以改回去
	CategoryId    uint           `gorm:"column:category_id;type:int(11);" json:"category_id"`
	ProblemId     uint           `gorm:"column:problem_id;type:int(11);" json:"problem_id"`
	CategoryBasic *CategoryBasic `gorm:"foreignKey:id;references:category_id"` //关联分类基础信息表
}

func (table *ProblemCategory) TableName() string {
	return "problem_category"
}
