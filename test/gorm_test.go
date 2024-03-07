package test

import (
	"gin_gorm_o/models"
	"testing"
)

/*
单元测试需要以_test.go结尾，goland才能识别出来
同时测试函数需要以Test开头
*/

func TestGormTest(t *testing.T) {
	//dsn := "root:123456@tcp(127.0.0.1:3306)/gin_gorm_oj?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	t.Fatal(err)
	//}

	//problem := &models.Problem{
	//	Identity:   "3",
	//	CategoryId: "3",
	//	Title:      "这是我的第三个问题",
	//	Content:    "请问3+3等于几",
	//	MaxMem:     100,
	//	MaxRuntime: 100,
	//}
	//db.Create(&problem)

	//data := make([]*models.Problem, 0)
	//err = db.Find(&data).Error
	//if err != nil {
	//	t.Fatal(err)
	//}
	//for _, v := range data {
	//	fmt.Println(v)
	//}

	err := models.DB.Debug().Model(new(models.CategoryBasic)).Where("name = ?", "数组").Delete(new(models.CategoryBasic)).Error

	if err != nil {
		t.Fatal(err)
	}
}
