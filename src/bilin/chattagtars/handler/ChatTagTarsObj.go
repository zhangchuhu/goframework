package handler

type ChatTagTarsObj struct {
}

func NewChatTagTarsObj() *ChatTagTarsObj {
	return &ChatTagTarsObj{}
}

const (
	MetricCodeSuccess      = iota //成功
	MetricCodeCreateErr           // 创建失败
	MetricCodeReadErr             // 读失败
	MetricCodeUpdateErr           // 更新失败
	MetricCodeDelErr              // 删除失败
	MetricCodeNoUserSetErr        // 无用户id

	MetricCodeDelAllWarn // 删除全部警告
)
