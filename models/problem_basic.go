package models

import (
	"gorm.io/gorm"
)

type ProblemBasic struct {
	gorm.Model
	Identity          string             `gorm:"column:identity;type:varchar(100);" json:"identity"` //问题的唯一标识
	Title             string             `gorm:"column:title;type:varchar(100);" json:"title"`
	Content           string             `gorm:"column:content;type:string;" json:"content"`
	MaxRuntime        int                `gorm:"column:max_runtime;type:int;" json:"max_runtime"`
	MaxMem            int                `gorm:"column:max_mem;type:int;" json:"max_mem"`
	PassNum           int64              `gorm:"column:pass_num;type:int" json:"pass_num"`     //完成问题个数，越多排名越高
	SubmitNum         int64              `gorm:"column:submit_num;type:int" json:"submit_num"` //提交次数，完成数相同时提交次数越少排名越高
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id;" json:"problem_categories"`
	TestCases         []*TestCase        `gorm:"foreignKey:problem_identity;references:identity;" json:"test_cases" ` //问题与测试用例关联表
}

func (table *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	//根据keyword进行模糊查询
	tx := DB.Model(new(ProblemBasic)).
		Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		Where("title like ? OR content like ?", "%"+keyword+"%", "%"+keyword+"%")
	if categoryIdentity != "" {
		tx.Joins("right join problem_category pc on pc.problem_id = problem_basic.id").
			Where("pc.category_id = (select cb.id from category_basic cb where cb.identity = ?)", categoryIdentity)
	}
	return tx
	//data := make([]*Problem, 0)
	//DB.Find(&data)
	//for _, v := range data {
	//	fmt.Println(v)
	//}
}
