package entity

type ClickInfo struct {
	UniqueID      string `json:"UniqueID"`
	Clicktime     string `json:"Clicktime"`
	IP            string `json:"IP"`
	OS            string `json:"OS"`
	Devicetype    string `json:"Devicetype"`
	Imei_md5      string `json:"Imei_md5"`
	IDFA          string `json:"IDFA"`
	MAC_MD5       string `json:"MAC_MD5"`
	Callback_url  string `json:"Callback_url"`
	From          string `json:"From"`
	Storagetime   string `json:"Storagetime"`
	Callback_time string `json:"Callback_time"`
	Callback_rsp  string `json:"Callback_rsp"`
}
