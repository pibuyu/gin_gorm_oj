package models

import "gorm.io/gorm"

type SubmitBasic struct {
	gorm.Model
	Identity        string        `gorm:"column:identity;type:varchar(100);" json:"identity"`
	UserIdentity    string        `gorm:"column:user_identity;type:varchar(100);" json:"user_identity"`
	ProblemBasic    *ProblemBasic `gorm:"foreignKey:identity;references:problem_identity"'` //关联问题基础表
	UserBasic       *UserBasic    `gorm:"foreignKey:identity;references:user_identity"'`    //关联用户基础表
	ProblemIdentity string        `gorm:"column:problem_identity;type:varchar(100);" json:"problem_identity"`
	Path            string        `gorm:"column:path;type:varchar(100);" json:"path"`
	Status          string        `gorm:"column:status;type:tinyint;" json:"status"` //-1-待判断；1-答案正确；2-答案错误；3-运行超时；4-运行超内存
}

func (table *SubmitBasic) TableName() string {
	return "submit_basic"
}

func GetSubmitList(problemIdentity, userIdentity string, status int) *gorm.DB {
	/*
		问题的content描述可能会比较长，我们把这个字段忽略掉(使content字段查询结果为空字符串)。
		又因为content是连接了ProblemBasic才获得的，所以在preload里处理
	*/
	tx := DB.Model(new(SubmitBasic)).
		Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB { return db.Omit("content") }).
		Preload("UserBasic")

	if problemIdentity != "" {
		tx.Where("problem_identity = ?", problemIdentity)
	}
	if userIdentity != "" {
		tx.Where("user_identity = ?", userIdentity)
	}
	if status != 0 {
		tx.Where("status = ?", status)
	}
	return tx
}
