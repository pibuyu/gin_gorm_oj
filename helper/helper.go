package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"gin_gorm_o/consts"
	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"net/smtp"
	"os"
	"time"
)

func SendVerifyCode(toUserEmail, content string) error {
	e := email.NewEmail()
	e.From = "<3531095171@qq.com>"
	e.To = []string{toUserEmail}
	e.Subject = "测试Golang邮件发送"
	e.HTML = []byte("<b>" + content + "</b>")
	//返回EOF的时候，关闭SSL重试
	return e.SendWithTLS(
		"smtp.qq.com:465",
		smtp.PlainAuth("",
			"3531095171@qq.com",
			"eyhbritymwqgcjca",
			"smtp.qq.com"),
		&tls.Config{
			ServerName:         "smtp.qq.com",
			InsecureSkipVerify: true,
		},
	)
}

// Get UUID
func GetUUID() string {
	return uuid.NewV4().String()
}

// MD5
func MD5(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

type UserClaims struct {
	Name     string `json:"name"`
	Identity string `json:"identity"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var secretKey = []byte(consts.SECRET_KEY)

// GenerateToken
// 生成token
func GenerateToken(name, identity string, is_admin int) (string, error) {
	//设置token为7天过期
	expireToken := time.Now().Add(time.Hour * 24 * 7).Unix()
	userClaim := &UserClaims{
		Identity: identity,
		Name:     name,
		IsAdmin:  is_admin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(expireToken),
		},
	}
	//生成加密token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	//生成盐加密的token
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// 解析token
func ParseToken(tokenString string) (*UserClaims, error) {
	userclaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userclaim,
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, err
	}

	return userclaim, nil
}

// CodeSave
func CodeSave(code []byte) (string, error) {
	dirName := "code/" + GetUUID()
	path := dirName + "/main.go"

	//创建保存文件夹
	err := os.Mkdir(dirName, 0777) // 0777为文件夹权限
	if err != nil {
		return "", err
	}

	//创建文件并把code写进文件
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	file.Write(code)
	defer file.Close()

	return path, nil
}
