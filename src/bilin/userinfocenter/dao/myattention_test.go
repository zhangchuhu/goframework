package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var err error
	AttentionDB, err = gorm.Open("mysql", "hujiao@HujiaoAttention:phU8o3l143@tcp(221.228.79.244:8066)/HujiaoAttention?readTimeout=500ms&charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("open db err", err)
		return
	}
	os.Exit(m.Run())
}

func TestMyAttention(t *testing.T) {
	count, err := MyAttentionNum(29481968)
	if err != nil {
		t.Error(err)
	}
	t.Log(count)
}
