include "common.thrift"

namespace java com.bilin.thriftserver.service.stub.gen
namespace go openstatus

/**
* 获取 用户审查状态
*/
service OpenStatusService extends common.BaseService{

    /**
    * 获取 用户审查状态
    * userId:用户uid
    *
    * return: 0: 未审核通过状态 1:审核通过
    *
    **/
	i32 getOpenStatusNew(1: i64 userId, 2: string version,3: string clientType,4: string ip )
    
    /**
    * 获取 用户审查状态
    * userId:用户uid
    *
    * return: 0: 未审核通过状态 1:审核通过
    *
    **/
	i32 getOpenStatus(1: i64 userId)

}