package config

import (
	"testing"
)

func TestInitAndSubConfig(t *testing.T) {
	if err := InitAndSubConfig("whatapp");err != nil{
		t.Error("InitAndSubConfig failed",err)
	}
}
