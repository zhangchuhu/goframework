namespace java com.bilin.thriftserver.service.stub.gen

/**
* 共通调用返回类型
*/
struct ComRet {
	1: required string result       //返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional string errorMsg		//如果result != "success"时，显示错误信息
}

/**
* 所有Service的基类Service
*/
service BaseService {
	i32 ping()
}