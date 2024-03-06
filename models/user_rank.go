package models

type UserRank struct {
	//gorm.Model
	Name             string `gorm:"column:name;type:varchar(100);" json:"name"`
	FinishProblemNum int64  `gorm:"column:finish_problem_num;type:int" json:"finish_problem_num"` //完成问题个数，越多排名越高
	SubmitNum        int64  `gorm:"column:submit_num;type:int" json:"submit_num"`                 //提交次数，完成数相同时提交次数越少排名越高
}

func (table *UserRank) TableName() string {
	return "user_rank"
}
