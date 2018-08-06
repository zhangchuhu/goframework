package dao_test

import (
	"bilin/userinfocenter/dao"
	_ "code.yy.com/yytars/goframework/tars/servant"
	"testing"
)

func TestGetUserAvatarInfo(t *testing.T) {
	info, err := dao.GetUserAvatarInfo(221414049, 17794899)
	if err != nil {
		t.Error("GetUserAvatarInfo error:" + err.Error())
	} else {
		t.Logf("GetUserAvatarInfo success, info: %v", info)
	}
}

func TestGetUserAvatatrInfos(t *testing.T) {
	info, err := dao.GetUserAvatatrInfos(201379, 100)
	if err != nil {
		t.Error("GetUserAvatatrInfos error:" + err.Error())
	} else {
		t.Logf("GetUserAvatatrInfos success, info: %v", info)
	}
}

func TestGetHttpsAvatar(t *testing.T) {
	url := dao.HttpsAvatarURL("http://img2.hujiaozhuanyi.com/imgs/201108/defaultBoy.png")
	if url != "https://img.inbilin.com/defaultBoy.png" {
		t.Error("HttpsAvatarURL failed")
	}
	t.Log(url)

	url = dao.HttpsAvatarURL("http://img.onbilin.com/30579309/30579309_1527012080689.jpg-small")
	if url != "https://img.inbilin.com/30579309/30579309_1527012080689.jpg-small" {
		t.Error("HttpsAvatarURL failed")
	}
	t.Log(url)

	url = dao.HttpsAvatarURL("http://img2.hujiaozhuanyi.com/imgs/201108/defaultGirl.png")
	if url != "https://img.inbilin.com/defaultGirl.png" {
		t.Error("HttpsAvatarURL failed")
	}
	t.Log(url)
}
