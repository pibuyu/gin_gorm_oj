package test

import (
	"fmt"
	"gin_gorm_o/helper"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func generateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
func TestSendEmail(t *testing.T) {
	//e := email.NewEmail()
	//e.From = "<3531095171@qq.com>"
	//e.To = []string{"2868319101@qq.com"}
	//e.Subject = "测试Golang邮件发送"
	//e.HTML = []byte("<b>收到请回复yes</b>")
	////返回EOF的时候，关闭SSL重试
	//e.SendWithTLS("smtp.qq.com:465",
	//	smtp.PlainAuth("",
	//		"3531095171@qq.com",
	//		"eyhbritymwqgcjca",
	//		"smtp.qq.com"),
	//	&tls.Config{
	//		ServerName:         "smtp.qq.com",
	//		InsecureSkipVerify: true,
	//	},
	//)
	code := generateCode(6)
	err := helper.SendVerifyCode("2098245863@qq.com", code)
	if err != nil {
		t.Fatalf("send email err:" + err.Error())
	}
}
