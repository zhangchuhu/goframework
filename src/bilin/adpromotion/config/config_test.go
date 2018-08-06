package config

import (
	"os"
	"testing"
)

func TestInitAndSubConfig(t *testing.T) {
	if err := InitAndSubConfig("whatapp"); err != nil {
		t.Error("InitAndSubConfig failed", err)
	}
}

func TestLoadconfig(t *testing.T) {
	os.Getwd()
	if err := loadconfig("whatapp"); err != nil {
		t.Error("InitAndSubConfig failed", err)
	}
}
