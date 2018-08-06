package entity

//用户数据
type UserRelationInfo struct {
	UserID      uint64 `json:"uid"`
	Nick        string `json:"nick"`
	Avatar      string `json:"avatar"`
	RelationVal int64  `json:"relation_value"`
	Headgear    string `json:"headgear"`
	MedalUrl    string `json:"medal_url"`
	MedalText   string `json:"medal_text"`
}

//榜单数据
type RelationStatistics struct {
	AnchorInfo   *UserRelationInfo   `json:"anchor_info"`   //主播信息
	RelationList []*UserRelationInfo `json:"relation_list"` //各用户亲密度信息
	Start        int                 `json:"start"`
	NumFound     int                 `json:"numFound"`
}
