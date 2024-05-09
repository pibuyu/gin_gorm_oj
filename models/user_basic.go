package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model

	Identity  string `gorm:"column:identity;type:varchar(100);" json:"identity"`
	Name      string `gorm:"column:name;type:varchar(100);" json:"name"`
	Password  string `gorm:"column:password;type:varvhar(100);" json:"password"`
	Phone     string `gorm:"column:phone;type:varchar(100);" json:"phone"`
	Mail      string `gorm:"column:mail;type:varvhar(100);" json:"mail"`
	PassNum   int64  `gorm:"column:pass_num;type:int" json:"pass_num"`     //完成问题个数，越多排名越高
	SubmitNum int64  `gorm:"column:submit_num;type:int" json:"submit_num"` //提交次数，完成数相同时提交次数越少排名越高
	IsAdmin   int    `gorm:"column:is_admin;type:int" json:"is_admin"`     //0-不是管理员，1-是管理员
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

type UserList []UserBasic

func (userlist *UserList) GetUserList() error {
	err := DB.Find(&userlist).Error
	return err
}
