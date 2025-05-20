package test

import (
	"GOIM/helper"
	"fmt"
	"testing"
)

func TestGenerateAndParseToken(t *testing.T) {
	id := "123"
	email := "123456"
	
	token, err := helper.GenerateToken(id, email)
	if err != nil {
		t.Errorf("generate token err = %v", err)
	}
	if token == "" {
		t.Error("generate token returned empty token")
	}
	fmt.Println(token)
	claims, err := helper.ParseToken(token)
	if err != nil {
		t.Errorf("parse token err = %v", err)
	}

	// 验证claims
	if claims.ID != id {
		t.Errorf("claims.ID = %v, id = %v", claims.ID, id)
	}
	if claims.Email != email {
		t.Errorf("claims.Email = %v, email = %v", claims.Email, email)
	}

	fmt.Println(claims)
}

func TestInvalidToken(t *testing.T) {
	// 测试无效token
	_, err := helper.ParseToken("invalid-token")
	if err == nil {
		t.Error("parse token should return error for invalid token")
	}
	fmt.Println(err)
}
