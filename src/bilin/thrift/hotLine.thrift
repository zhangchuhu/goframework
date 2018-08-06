include "common.thrift"

namespace java com.bilin.thriftserver.service.stub.gen

/**
* JoinHotLine接口调用返回接口
*/
struct JoinHotLineRet {
	1: required string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional	string	errorMsg			//如果result != success，则返回错误消息
	3: required i32	status					//以如下方式进入成功 1:管理员 2：主播 3:听众  以如下方式进入失败: -1:权限限制（含反垃圾限制) -2:直播异常（未开始或已结束或不存在）
	4: optional	string	nickname			//用户昵称
	5: optional	string	headerUrl			//用户头像
	6: required	i32	sex						//性别 0=女 1=男
	7: required	i32	age						//年龄
	8: optional	string	cityName			//城市
    9: optional i64 hostBilinId         //主播比邻ID
    10: optional string title           //直播标题
    11: optional i32 roomType           //房间类型 1、官方频道；2、PGC；3、UGC
    12: optional i32 roomCategoryId         //直播分类id
    13: optional i32 roomPendantLevel           //主播挂件等级
    14: optional string sign			//用户sign
}

/**
* GetInitHotLineInfo
*/
struct GetInitHotLineInfoRet {
	1: required string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional	string	errorMsg			//如果result != success，则返回错误消息
	3: required double	viewUsersCoff		//观看人数系数
	4: required double	praisesCoff			//赞数系数
	5: required i64	viewCnt					//总观看人次(真实)
	6: required i64	praises					//总赞数(真实)
	7: optional i32 lineType                //直播类型0:热线直播 1:视频直播，默认0
        8: optional i64	currentViewCntIncr			//本次直播新增的总观看人次(真实)
}

/**
* GetHotLineNoticeText接口调用返回接口
*/
struct GetHotLineNoticeTextRet {
	1: required string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional	string	errorMsg			//如果result != success，则返回错误消息
	3: required string	noticeText			//系统公告
}

enum TokenAuthType {
    tat_up_voice = 1,
    tat_down_voice = 2,
    tat_up_video = 3,
    tat_down_video = 4,
    tat_up_text = 5,
    tat_down_text = 6,
}

/**
*热线直播相关服务接口
*/
service HotLineService extends common.BaseService {

	//刷数据
	//hotLineId			热线ID
	//cuu	    		当时在线观看人数(系数)
	//tvc				总观看次数(系数)
	//ps				赞总数(系数)
	//tcuu	    		当时在线观看人数(真实)
	//ttvc				总观看次数(真实)
	//tps				赞总数(真实)
	common.ComRet freshData(1: i32 hotLineId, 2: i64 cvu, 3: i64 tvc, 4: i64 ps, 5: i64 tcvu, 6: i64 ttvc, 7: i64 tps)

	//用户进入热线直播
	//hotLineId			热线ID
	//userId			用户ID
	JoinHotLineRet joinHotLine(1: i32 hotLineId, 2: i64 userId)	

	//获取热线直播初期数据
	GetInitHotLineInfoRet getInitHotLineInfo(1: i32 hotLineId)
	
	//踢人
	//hotLineId 	热线ID
	//userId    	踢人用户ID
	//targetUserId	被踢用户ID
	//targetType	被题用户类型(1:管理员 2：主播 3:听众)
	common.ComRet kickUser(1: i32 hotLineId, 2: i64 userId, 3: i64 targetUserId, 4: i32 targetType)
	
	//获取热线直播公告
	GetHotLineNoticeTextRet GetHotLineNoticeText(1: i32 hotLineId);
	
	//主播离开
	common.ComRet hostUserLeave(1: i32 hotLineId, 2: i64 userId);
	
	//主播太长时间不在线
	common.ComRet hostUserOfflineTooLong(1: i32 hotLineId, 2: i64 userId);

	//更新直播间内哪些人有哪些权限
	//liveId 直播id，热线为热线直播id（hotLineId），视频为视频直播id（videoLiveId）
	//userList uid列表
	//authType 权限type。
	//disable =true 为禁用对应权限，=false为开启对应权限
	common.ComRet updateLiveAuthority(1: i32 liveId, 2: list<i64> userList, 3: list<TokenAuthType> authType, 4: bool disable);
}


struct TaskReq {
    1:i64 uid 
    2:string task_key
    3:string task_id
    4:i64 room_id
}

service DataService extends common.BaseService
{
    i32 start(1:list<TaskReq> tasks);
    i32 cancel(1:list<TaskReq> tasks);
}
