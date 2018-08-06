package sensitive

import (
	"fmt"
	"testing"
)

func init() {
	if err := LoadSensiTive(); err != nil {
		fmt.Println("LoadSensiTive failed", err)
	}
}
func TestSensiTiveWord(t *testing.T) {
	if !SensiTiveWord("我们来寂寞") {
		t.Error("should be sensitive")
	}
}
