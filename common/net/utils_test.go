package net

import (
	"fmt"
	"testing"
)

func TestGetIP(t *testing.T) {
	ip, err := GetIP()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ip)
}
