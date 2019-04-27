package flow

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var secretKey = []byte("abcd1234")

//定义
type PermisstionClaim struct{
	UserId string `json:"userId"`
	Name string `json:"name"`
	RoleId string `json:"roleId"`
	RoleName string  `json:"roleName"`

	jwt.StandardClaims
}

//jwKeyFunc返回密钥
func jwtKeyFunc(token *jwt.Token)(interface{},error){
	return secretKey,nil
}

//sign 生成token
func Sign(name,uid,roleId,roleName string)(string,error){
	fmt.Println(name,uid,roleId,roleName )
	//设置过期时间  演示设置2 分钟
	expAt := time.Now().Add(time.Duration(2)*time.Minute).Unix()

	//创建声明
	claims := PermisstionClaim{
		UserId: uid,
		Name:   name,
		RoleId: roleId,
		RoleName:roleName,
		StandardClaims : jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    "system",
		},
	}

	//创建token，指定加密算法为HS256
	token := jwt.NewWithClaims(jwt.SigningMethodES256,claims)
	fmt.Println(token)
	//生成token 有错误
	//return token.SignedString(secretKey)

	//不带key 可以生成TOKEN
	return token.SigningString()
}