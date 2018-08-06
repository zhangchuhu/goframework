package handler

import (
	"fmt"
	"testing"
)

func TestIsWhite(t *testing.T) {
	fmt.Println("TestIsWhite")
	IsWhite(1, 1)
	IsWhite(1, 0)
}

func TestIsWhite01(t *testing.T) {
	fmt.Println("TestIsWhite01")
	AddFemaleWhite(1)
	fmt.Println("FemaleWhite")
	fmt.Println(IsWhite(1, 1))
	DelFemaleWhite(1)
	fmt.Println("FemaleWhite")
	fmt.Println(IsWhite(1, 1))
	fmt.Println("")

	AddMaleWhite(1)
	AddMaleWhite(1)
	AddMaleWhite(1)
	AddMaleWhite(1)
	AddMaleWhite(1)
	AddMaleWhite(1)
	fmt.Println("MaleWhite")
	fmt.Println(IsWhite(1, 0))
	DelMaleWhite(1)
	fmt.Println("MaleWhite")
	fmt.Println(IsWhite(1, 0))
}
