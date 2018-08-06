include "common.thrift"

namespace java com.bilin.thriftserver.service.stub.gen

/**
* getUserRoom接口调用返回接口
*/
struct GetUserRoomRet {
	1: required string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional	string	errorMsg			//如果result != success，则返回错误消息
	3: required bool	existed				//是否存在
	4: optional string	name				//房间名称
	5: optional string	topic				//房间话题
}

struct QueryMyRoomWhiteListRet {
	1: required string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional	string	errorMsg			//如果result != success，则返回错误消息
	3: optional list<i32> whiteList
}

/**
*呼叫记录相关服务接口
*/
service CallRecordService extends common.BaseService {

	//完成新用户随机呼引导任务接口
	//toUserId			用户ID
	//isFinish			是否听完录音
	//flowerNum			用户接收花朵数
	common.ComRet finishNewUserRandomCallTask(1: i64 toUserId, 2: bool isFinish, 3: i32 flowerNum)	

	//CC服务添加未接来电记录
	//fromUserId		主叫方用户ID
	//toUserId			被叫方用户ID
	//isFriendCall		是否好友直呼 （1：好友直呼，2：选择呼）
	common.ComRet addMissedCall(1: i64 fromUserId, 2: i64 toUserId, 3: i32 isFriendCall, 4: string applyId)

	//CC服务方提供的通话记录入库
	//beginTime			开始时间（1970年到现在的毫秒数）
	//callId			通话ID
	//endTime			结束时间（1970年到现在的毫秒数）
	//flowerCounts		花朵数：送花用户-接收用户-数量,送花用户-接收用户-数量
	//toUserId			被叫方用户ID
	//callType			通话类型1：好友直呼，2：随机呼叫，3：选择呼叫
	//fromUserId		主叫方用户ID
	//netType			网络类型
	common.ComRet addCallRecordByCCServer(1: i64 beginTime, 2: string callId, 3: i64 endTime, 4: string flowerCounts, 5: i64 toUserId, 6: i32 callType, 7: i64 fromUserId, 8: string netType)

	//CC服务送花
	//callId			通话ID
	//fromUserId		送花用户ID
	//toUserId			接收用户ID
	//type				来源类型（1：CC系统，2：手机客户端）
	common.ComRet addUserFlowerByCall(1: string callId, 2: i64 fromUserId, 3: i64 toUserId, 4: i32 type)
	
	//获取用户房间
	//userId			用户ID
	//roomId			房间ID
	//roomType			房间类型
	GetUserRoomRet getUserRoom(1: i64 userId, 2: i64 roomId, 3: i32 roomType)
	
	//修改用户房间话题
	//userId			用户ID
	//roomId			房间ID
	//roomType			房间类型
	//newTopic			新房间话题
	common.ComRet updateRoomTopic(1: i64 userId, 2: i64 roomId, 3: i32 roomType, 4: string newTopic)
	
	//查询自建房间白名单
	QueryMyRoomWhiteListRet queryMyRoomWhiteList();
}


