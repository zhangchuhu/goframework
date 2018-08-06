package userinfo

import (
	"bilin/clientcenter"
	"bilin/protocol/userinfocenter"
	"context"
	"time"
)

func IsAppleCheckUser(uid uint64, version, clienttype, ip string) bool {
	if clienttype != "IPHONE" {
		return false
	}
	ctx, cancle := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancle()
	resp, err := clientcenter.UserInfoClient().IsAppleCheckUser(ctx, &userinfocenter.IsAppleCheckUserReq{
		Uid:        uid,
		Version:    version,
		Clienttype: clienttype,
		Ip:         ip,
	})
	if err != nil || resp == nil {
		return false
	}
	return resp.Applecheckuser
}
