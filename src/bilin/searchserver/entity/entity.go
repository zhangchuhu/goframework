package entity

import (
	u "bilin/searchserver/updater"
)

type UserE struct {
	u.UserU
	Avatar      string `json:"avatar"`
	Sex         string `json:"sex"`
	Age         string `json:"age"`
	Location    string `json:"location"`
	Live        string `json:"live"`
	RoomUserNum string `json:"room_user_num"`
}

type RoomE struct {
	u.RoomU
	Avatar  string   `json:"avatar"`
	StartAt string   `json:"start_at"`
	UserNum string   `json:"user_num"`
	TagURL  []string `json:"tag_url,omitempty"`
}

type SongE struct {
	u.SongU
	Duration        string `json:"duration"`
	UploadBy        string `json:"upload_by"`
	Lyric           string `json:"lyric"`
	LyricMd5        string `json:"lyric_md5"`
	LyricLen        string `json:"lyric_len"`
	Audio           string `json:"audio"`
	AudioMd5        string `json:"audio_md5"`
	AudioLen        string `json:"audio_len"`
	Instrumental    string `json:"instrumental"`
	InstrumentalMd5 string `json:"instrumental_md5"`
	InstrumentalLen string `json:"instrumental_len"`
	Pkg             string `json:"pkg"`
	PkgMd5          string `json:"pkg_md5"`
	PkgLen          string `json:"pkg_len"`
}
