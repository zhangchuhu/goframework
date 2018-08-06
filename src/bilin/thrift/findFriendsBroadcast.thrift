include "common.thrift"

 namespace java com.bilin.thriftserver.service.stub.gen

 struct OfflineFindFriendsBroadcastRet {
 	1: required	string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
 	2: optional string	errorMsg			//如果result != success，则返回错误消息
 }


 service FindFriendsBroadcastService extends common.BaseService {

     //userIdList:用户id列表
 	OfflineFindFriendsBroadcastRet offlineFindFriendsBroadcastByUserIdList(1: list<i64> userIdList)

 }