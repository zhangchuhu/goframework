include "common.thrift"

namespace java com.bilin.thriftserver.service.stub.gen

struct OfficialHotlineRet {
	1: required	string	result				//返回结果是否成功,如果为"success"  则调用成功，否则调用失败
	2: optional	string	response			//返回结果json格式
	3: optional string	errorMsg			//如果result != success，则返回错误消息
}

/**
* 官频接口接口
*/
service OfficialHotlineService extends common.BaseService {

    /**
    * 官频切换回调
    * sid 频道号
    * oldUserId 原主播
    * newUserId 新主播
    * msg 业务方的透传消息
    * resultCode -1 用户不在房间
		0 成功
		-501 查询房间主播身份（join_hot_line）thrift请求异常
		-502 查询房间主播身份（join_hot_line）身份失败，身份值小于0
    */
    OfficialHotlineRet switchCallback(1: i32 sid, 2: i64 oldUserId, 3: i64 newUserId, 4: string msg, 5: i32 resultCode)
    /**
     * 官频主播上麦成功调用
     * sid 频道号
     * userId 主播ID
     */
    OfficialHotlineRet onOfficialMike(1: i32 sid, 2: i64 userId);

}
