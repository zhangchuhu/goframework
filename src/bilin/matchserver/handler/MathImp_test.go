package handler

import (
	"fmt"
	"testing"
	"time"
)

var MatchSexMap map[string][]string
var MatchNoSexList [][]string

// newUser(uid uint32, matchType int, sex int, role int, province string)
func TestNewUser(t *testing.T) {
	newUser(1, 0, 0, 1, "gd") // 异性，男，白
	newUser(2, 0, 1, 1, "gd") // 异性，女，白
	newUser(3, 0, 0, 0, "gd") // 异性，男，普
	newUser(4, 0, 1, 0, "gd") // 异性，女，普
	newUser(5, 1, 1, 0, "gd") // 同性，女，普
	newUser(6, 1, 0, 0, "gd") // 同性，男，普
	newUser(7, 1, 0, 1, "gd") // 同性，男，白
	ClearRedisData()
}

func TestHandleMatchSex01(t *testing.T) {
	fmt.Println("")
	fmt.Println("TestHandleMatchSex01")

	user1, _ := newUser(1, 0, 0, 1, "gd") // 异性，男，白
	user2, _ := newUser(2, 0, 0, 0, "gd") // 异性，男，普
	user3, _ := newUser(3, 0, 0, 1, "gd") // 异性，男，白
	user4, _ := newUser(4, 0, 1, 0, "gd") // 异性，女，普

	HandleMatchSex()

	userlist := MatchSexMap[user4]
	fmt.Println("Match OK", user4, "##", userlist)

	if userlist[0] == user1 && userlist[1] == user3 && userlist[2] == user2 {
		fmt.Println("TestHandleMatchSex01 ok")
	} else {
		fmt.Println("TestHandleMatchSex01 no ok")
	}

	ClearRedisData()
}

func TestHandleMatchSex02(t *testing.T) {
	fmt.Println("")
	fmt.Println("TestHandleMatchSex02")

	user1, _ := newUser(1, 0, 0, 1, "gd") // 异性，男，白
	newUser(2, 0, 0, 0, "gd")             // 异性，男，普
	newUser(3, 0, 0, 0, "gd")             // 异性，男，普
	user4, _ := newUser(4, 0, 1, 1, "gd") // 异性，女，白
	fmt.Println("")
	HandleMatchSex()

	userlist := MatchSexMap[user4]
	fmt.Println("Match OK", user4, "##", userlist)

	if userlist[0] == user1 {
		fmt.Println("TestHandleMatchSex02 ok")
	} else {
		fmt.Println("TestHandleMatchSex02 no ok")
	}

	ClearRedisData()
}

func TestHandleMatchSex03(t *testing.T) {
	fmt.Println("")
	fmt.Println("TestHandleMatchSex03")

	user1, _ := newUser(1, 0, 0, 1, "gd") // 异性，男，白
	newUser(2, 0, 0, 0, "gd")             // 异性，男，普
	user3, _ := newUser(3, 0, 0, 1, "gd") // 异性，男，白
	user4, _ := newUser(4, 0, 1, 1, "gd") // 异性，女，白
	user5, _ := newUser(5, 0, 1, 1, "gd") // 异性，女，白
	newUser(6, 0, 1, 1, "gd")             // 异性，女，白
	fmt.Println("")
	HandleMatchSex()

	userlist := MatchSexMap[user4]
	fmt.Println("Match OK", user4, "##", userlist)

	if userlist[0] == user1 {
		fmt.Println("TestHandleMatchSex03 ok")
	} else {
		fmt.Println("TestHandleMatchSex03 no ok")
	}

	userlist = MatchSexMap[user5]
	fmt.Println("Match OK", user5, "##", userlist)

	if userlist[0] == user3 {
		fmt.Println("TestHandleMatchSex03 ok")
	} else {
		fmt.Println("TestHandleMatchSex03 no ok")
	}

	ClearRedisData()
}

func TestHandleMatchSex04(t *testing.T) {
	fmt.Println("")
	fmt.Println("TestHandleMatchSex04")

	user1, _ := newUser(1, 0, 0, 1, "gd")   // 异性，男，白
	user2, _ := newUser(2, 0, 0, 0, "gd")   // 异性，男，普
	user3, _ := newUser(3, 0, 0, 1, "gd")   // 异性，男，白
	user4, _ := newUser(4, 0, 0, 0, "gd")   // 异性，男，普
	user5, _ := newUser(5, 0, 0, 0, "gd")   // 异性，男，普
	user6, _ := newUser(6, 0, 0, 0, "gd")   // 异性，男，普
	user7, _ := newUser(7, 0, 0, 0, "gd")   // 异性，男，普
	user8, _ := newUser(8, 0, 1, 0, "gd")   // 异性，女，普
	user9, _ := newUser(9, 0, 1, 1, "gd")   // 异性，女，白
	user10, _ := newUser(10, 0, 1, 0, "gd") // 异性，女，普
	fmt.Println("")
	HandleMatchSex()

	userlist := MatchSexMap[user9]
	fmt.Println("Match OK", user9, "##", userlist)

	if userlist[0] == user1 {
		fmt.Println("TestHandleMatchSex04 ok")
	} else {
		fmt.Println("TestHandleMatchSex04 no ok")
	}

	userlist = MatchSexMap[user8]
	fmt.Println("Match OK", user8, "##", userlist)

	if userlist[0] == user3 && userlist[1] == user2 && userlist[2] == user4 {
		fmt.Println("TestHandleMatchSex04 ok")
	} else {
		fmt.Println("TestHandleMatchSex04 no ok")
	}

	userlist = MatchSexMap[user10]
	fmt.Println("Match OK", user10, "##", userlist)

	if userlist[0] == user5 && userlist[1] == user6 && userlist[2] == user7 {
		fmt.Println("TestHandleMatchSex04 ok")
	} else {
		fmt.Println("TestHandleMatchSex04 no ok")
	}

	ClearRedisData()
}

func WhileMatch() {
	i := 0
	for i < 100 {
		HandleMatchSex()
		time.Sleep(10 * time.Millisecond)
		i++
	}
}
func TestHandleMatchSex05(t *testing.T) {
	fmt.Println("")
	fmt.Println("TestHandleMatchSex05")
	go WhileMatch()
	for i := 0; i < 100; i++ {
		newUser(uint32(i), 0, 0, 0, "gd") // 异性，男，普
	}

	for i := 100; i < 110; i++ {
		newUser(uint32(i), 0, 1, 1, "gd") // 异性，女，白
	}

	for i := 110; i < 133; i++ {
		newUser(uint32(i), 0, 1, 0, "gd") // 异性，女，普
	}

	fmt.Println("")
	HandleMatchSex()

	HandleMatchSex()

	if len(MatchSexMap) == 33 {
		fmt.Println("Match OK", len(MatchSexMap))
	} else {
		fmt.Println("Match no OK", len(MatchSexMap))
	}

	ClearRedisData()
}

func TestHandleMatchSex06(t *testing.T) {
	fmt.Println("")
	fmt.Println("TestHandleMatchSex06")

	for i := 0; i < 30; i++ {
		newUser(uint32(i), 0, 0, 0, "gd") // 异性，男，普
	}

	for i := 30; i < 60; i++ {
		newUser(uint32(i), 0, 0, 0, "gx") // 异性，男，普,gx
	}

	for i := 60; i < 70; i++ {
		newUser(uint32(i), 0, 1, 1, "gd") // 异性，女，白
	}

	for i := 70; i < 80; i++ {
		newUser(uint32(i), 0, 1, 0, "gd") // 异性，女，普
	}

	fmt.Println("")
	HandleMatchSex()

	HandleMatchSex()

	if len(MatchSexMap) == 20 {
		fmt.Println("Match OK", len(MatchSexMap))
	} else {
		fmt.Println("Match no OK", len(MatchSexMap))
	}

	ClearRedisData()
}

func TestHandleMatchNoSex01(t *testing.T) {
	fmt.Println("")
	fmt.Println("TestHandleMatchNoSex01")

	user1, _ := newUser(1, 1, 0, 1, "gd") // 同性，男，白
	user2, _ := newUser(2, 1, 0, 0, "gd") // 同性，男，普
	user3, _ := newUser(3, 1, 0, 1, "gd") // 同性，男，白
	user4, _ := newUser(4, 1, 0, 0, "gd") // 同性，男，普
	user5, _ := newUser(5, 1, 0, 0, "gd") // 同性，男，普
	user6, _ := newUser(6, 1, 0, 0, "gd") // 同性，男，普
	user7, _ := newUser(7, 1, 0, 0, "gd") // 同性，男，普
	user8, _ := newUser(8, 1, 0, 0, "gd") // 同性，男，普
	newUser(9, 1, 0, 0, "gd")             // 同性，男，普
	newUser(10, 1, 0, 0, "gd")            // 同性，男，普

	fmt.Println("")
	HandleMatchNoSex()

	userlist := MatchNoSexList[0]
	fmt.Println(userlist)
	if userlist[0] == user1 && userlist[1] == user3 && userlist[2] == user2 && userlist[3] == user4 {
		fmt.Println("TestHandleMatchSex04 ok")
	} else {
		fmt.Println("TestHandleMatchSex04 no ok")
	}

	userlist = MatchNoSexList[1]
	fmt.Println(userlist)
	if userlist[0] == user5 && userlist[1] == user6 && userlist[2] == user7 && userlist[3] == user8 {
		fmt.Println("TestHandleMatchSex04 ok")
	} else {
		fmt.Println("TestHandleMatchSex04 no ok")
	}
	ClearRedisData()
}

func TestHandleMatchNoSex02(t *testing.T) {
	fmt.Println("")
	fmt.Println("TestHandleMatchNoSex02")

	user1, _ := newUser(1, 1, 1, 0, "gd") // 同性，女，普
	user2, _ := newUser(2, 1, 1, 0, "gd") // 同性，女，普
	user3, _ := newUser(3, 1, 1, 0, "gd") // 同性，女，普
	user4, _ := newUser(4, 1, 1, 0, "gd") // 同性，女，普
	user5, _ := newUser(5, 1, 1, 0, "gd") // 同性，女，普
	user6, _ := newUser(6, 1, 1, 0, "gd") // 同性，女，普
	user7, _ := newUser(7, 1, 1, 0, "gd") // 同性，女，普
	user8, _ := newUser(8, 1, 1, 0, "gd") // 同性，女，普
	newUser(9, 1, 1, 0, "gd")             // 同性，女，普
	newUser(10, 1, 1, 0, "gd")            // 同性，女，普

	fmt.Println("")
	HandleMatchNoSex()

	userlist := MatchNoSexList[0]
	fmt.Println(userlist)
	if userlist[0] == user1 && userlist[1] == user2 && userlist[2] == user3 && userlist[3] == user4 {
		fmt.Println("TestHandleMatchNoSex02 ok")
	} else {
		fmt.Println("TestHandleMatchNoSex02 no ok")
	}

	userlist = MatchNoSexList[1]
	fmt.Println(userlist)
	if userlist[0] == user5 && userlist[1] == user6 && userlist[2] == user7 && userlist[3] == user8 {
		fmt.Println("TestHandleMatchNoSex02 ok")
	} else {
		fmt.Println("TestHandleMatchNoSex02 no ok")
	}

	ClearRedisData()
}
