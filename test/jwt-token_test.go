package test

import (
	"fmt"
	"gin_gorm_o/consts"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

type userClaims struct {
	Name     string `json:"name"`
	Identity string `json:"identity"`
	jwt.StandardClaims
}

var secretKey = []byte(consts.SECRET_KEY)

//生成token
/*
声明claims，要继承StandardClaims: &jwt.StandardClaims{},
调用jwt.NewWithClaims对claim加密，得到第二部分
再对token进行盐加密，得到第三部分
*/
func TestGenerateToken(t *testing.T) {
	userClaim := &userClaims{
		Identity:       "243045",
		Name:           "zhang",
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("sign token process error:" + err.Error())
	}
	fmt.Println(signedToken)
	fmt.Println(len(signedToken))
}

// 解析token
func TestParseToken(t *testing.T) {
	token_string := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiaGhmIiwiaWRlbnRpdHkiOiIxIn0.f2R7NFOfC9nLQ-TCB3j4gFU9FwTSw8sXQsi1moO6CsI"
	userclaim := new(userClaims)
	claims, err := jwt.ParseWithClaims(token_string, userclaim,
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
	if err != nil {
		t.Fatal(err)
	}
	if claims.Valid {
		fmt.Println(userclaim)
	}
}
