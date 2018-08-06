package entity

type HeadgearInfo struct {
	Uid        int64  `json:"userid"`
	Headgear   string `json:"headgear"`
	EffectTime string `json:"effecttime"`
	ExpireTime string `json:"expiretime"`
	Id         int64  `json:"id"`
}
