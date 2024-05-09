package test

import (
	"fmt"
	"gin_gorm_o/models"
	"testing"
)

func TestGetUserList(t *testing.T) {
	list := new(models.UserList)
	err := list.GetUserList()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(list)
}
