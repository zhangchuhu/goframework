# 敏感词服务代理
namespace cpp server.bilin

service MsgFilter {
    i32 check_msg(1: string msg)
}
