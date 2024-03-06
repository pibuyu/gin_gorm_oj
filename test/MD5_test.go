package test

import (
	"fmt"
	"gin_gorm_o/helper"
	"testing"
)

func TestMD5(t *testing.T) {
	md5 := helper.MD5("123456")
	fmt.Println(md5)
}
