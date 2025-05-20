package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	JWTSecret = "your-secret-key"
	TokenExpireDuration = 24*time.Hour
)

type UserClaims struct{
	ID bson.ObjectID `json:"_id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(ID bson.ObjectID, email string)(string, error) {
	claims := &UserClaims{
		ID: ID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	
	// 使用密钥签名
	return token.SignedString([]byte(JWTSecret))
}

// 解析JWT Token
func ParseToken(tokenString string)(*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token)(interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	// 验证token
	if claims, ok := token.Claims.(*UserClaims);ok&&token.Valid {
		return claims,nil
	}
	return nil, errors.New("invalid token")
}