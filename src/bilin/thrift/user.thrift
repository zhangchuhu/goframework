include "common.thrift"

namespace java com.bilin.thriftserver.service.stub.gen

/**
*通话过程中加关注调用返回类型
*/
struct GetUserForConServerRet {
	1: required	string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional string	errorMsg			//如果result != success，则返回错误消息
	3: optional	string	clientType			//客户端类型
}

/**
*查询关注用户的人数
*/
struct QueryAttentionMeCountRet {
	1: required	string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional string	errorMsg			//如果result != success，则返回错误消息
	3: optional map<i64,i64> attentionMeCountMap	//关注人数map  userId ---> count	
}


/**
*用户相关服务接口
*/
service UserService extends common.BaseService {
	//禁用用户
	//userIds 用户ID,多个用户id之间以","隔开
	//bilinIds 用户比邻ID,多个用户比邻id之间以","隔开
	common.ComRet internalHttpForbidUser(1: string userIds, 2: string bilinIds)

	//con_server服务端获取用户信息
	//userId				用户ID
	//fromUserId			用户ID
	//groupId				讨论组ID（=0：正常电话push，>0：讨论组召唤push）
	//ifPush				是否push
	GetUserForConServerRet getUserForConServer(1: i64 userId, 2: i64 fromUserId, 3: i64 groupId, 4: bool ifPush, 5: string requestType)	
	
	//查询用户的关注人数
	QueryAttentionMeCountRet QueryAttentionMeCount(1: list<i64> userList)
}