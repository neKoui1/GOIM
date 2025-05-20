package test

import (
	"GOIM/helper"
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	pwd := "testPWD"
	hash, err := helper.HashPassword(pwd)
	if err != nil {
		t.Fatal("hash password error", err)
	}
	fmt.Println("hash password", hash)

	fmt.Println("check password", helper.CheckPassword(pwd, hash))
}