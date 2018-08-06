package handler

const (
	SUCCESS = iota
	ParamInvalid
	ParamMarshalFailed
	GetDBUserInfoFailed
	GetCacheUserInfoFailed
	SetCacheUserInfoFailed
	DelCacheUserInfoFailed
	GetThriftOpenStatusFailed
	GetCacheOpenStatusFailed
	SetCacheOpenStatusFailed
	BatchUserBLNumFailed
	AttentionMeFailed
)

const (
	GetDBUserInfoFailedKey               = "get_userinfo_db_fail"
	GetCacheUserInfoFailedKey            = "get_userinfo_cache_fail"
	SetCacheUserInfoFailedKey            = "set_userinfo_cache_fail"
	DelCacheUserInfoFailedKey            = "del_userinfo_cache_fail"
	GetThriftOpenStatusFailedKey         = "get_userinfo_thrift_fail"
	GetCacheOpenStatusFailedKey          = "get_userinfo_openstatus_cache_fail"
	SetCacheOpenStatusFailedKey          = "set_userinfo_openstatus_cache_fail"
	GetDBUserInfoTableUserFailedKey      = "set_userinfo_table_user_fail"
	GetDBUserInfoTableAvatarFailedKey    = "set_userinfo_table_avatar_fail"
	GetDBUserInfoTableExtensionFailedKey = "get_userinfo_table_extension_fail"
	GetDBUserInfoTableCountryFailedKey   = "get_userinfo_table_country_fail"
	UpdateUserInfoMarshalFaileKey        = "update_userinfo_marshal_fail"
	GetDBUserInfoKey                     = "get_db_userinfo"
	GetCacheUserInfoKey                  = "get_cahce_userinfo"
	GetThriftOpenStatusKey               = "get_userinfo_thrift"
	GetCacheOpenStatusKey                = "get_userinfo_openstatus_cache"
)
