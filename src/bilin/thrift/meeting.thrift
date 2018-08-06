include "common.thrift"
namespace java com.bilin.thriftserver.service.stub.gen

/**
*会议话题
*/
struct Topic {
	1: required	i64		id		//会议话题ID
	2: required string	topic	//会议话题
}

/**
*CC服务方需要的获取所有的会议话题的配置
*/
struct QueryAllMeetingTopicListByCCServerRet {
	1: required string	result					//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional string	errorMsg				//如果result != success，则返回错误消息
	3: optional list<Topic>	meetingTopicList	//会议话题列表
}

/**
*会议昵称
*/
struct Nickname {
	1: required	i64		id					//
	2: required	list<string>	boyNick		//男性可用昵称列表
	3: required	list<string>	girlNick	//女性可用昵称列表
	4: required	string	title				//标题
}

/**
*获取所有的会议昵称配置
*/
struct QueryAllMeetingNicknameByCCServerRet {
	1: required	string	result					//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional string	errorMsg				//如果result != success，则返回错误消息
	3: optional	list<Nickname>	nicknameList	//加完关注后，两用户之间的关注关系   0: 两个互不关注, 2: 两个互相关注 , 1 : A 关注了 B, -1: B关注了A
}


/**
*用户设备信息
*/
struct UserAgent {
	1: required	i64		userId		//用户ID
	2: required	string	version		//客户端版本
	3: required	string	clientType	//客户端类型
	4: required	string	osVersion	//客户端操作系统版本
}

/**
*批量获取用户设备信息接口调用返回对象
*/
struct QueryUserDeviceInfoByCCServerRet {
	1: required string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional	string	errorMsg			//如果result != success，则返回错误消息
	3: optional	list<UserAgent> userList		//用户设备列表	
}

/**
*根据话题类型获取话题列表接口调用返回对象
*/
struct QueryTopicListByCCServerRet {
	1: required	string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional	string	errorMsg			//如果result != success，则返回错误消息
	3: optional	list<Topic> meettingTopicList	//用户设备列表	
}

/**
*超能力标签
*/
struct SuperPowerTag {
	1: required	i64		id					//超能力标签ID
	2: required	string	name				//超能力标签名称
	3: required	string	iconImgUrl			//超能力标签Icon图片地址
	4: required	string	tagImgUrl			//超能力标签图片地址
}

/**
*获取所有的会议超能力配置接口调用返回对象
*/
struct QueryAllSuperPowerTagListByCCServerRet {
	1: required	string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional	string	errorMsg			//如果result != success，则返回错误消息
	3: optional	list<SuperPowerTag> tagList	//超能力列表	
}
/**
*是否能踢人接口调用返回对象
*/
struct IsCanKickOutRet {
	1: required	string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional	string	errorMsg			//如果result != success，则返回错误消息
	3: optional i32 	isCanKickout		//是否能踢人 1:可以踢，0:不能踢
	4: optional i32 	isNewPipe			//1:新通道，0：旧通道
	5: optional i32 	forbidLevel			//第1,2,3级规则，3表示把房主踢出
	6: optional string	forbidPrompt		//禁止踢人文案，如果不能踢人（isCanKickout=false）时返回
}

/**
*用户信息
*/
struct MeetingUser {
	1: required  i64 	userId				//用户ID
	2: optional  string	nickname			//昵称
	3: optional  string	rcUrl				//随机呼头像
	4: optional	 i32 	sex					//性别
	5: optional	 i32	age					//年龄
	6: optional  i64    numOfFlower			//花朵数
	7: optional  string smallUrl			//小头像
}
/**
*批量获取会议用户信息调用返回对象
*/
struct QueryUserListByMeetingRet {
	1: required	string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional	string	errorMsg			//如果result != success，则返回错误消息
	3: optional	list<MeetingUser> userList	//群聊中所有用户信息	
}

struct IsUserInRandomCallingWhiteListRet {
	1: required	string	result				//返回结果,如果为"success"  则调用成功，否则调用失败
	2: optional	string	errorMsg			//如果result != success，则返回错误消息
	3: required	i32		ret					//是否在白名单中　０:不在　1:在
}
struct SpamUserLevel {
    1: i32 level;       ////用户等级 普通用户：3,  黑名单：2,  灰名单：1, 白名单：0
    2: bool isCheat;     //是否外挂用户
}


/**
*会议相关服务接口
*/
service MeetingService extends common.BaseService{
	//获取所有的会议话题配置
	QueryAllMeetingTopicListByCCServerRet queryAllMeetingTopicListByCCServer()

	//获取所有的会议昵称配置
	QueryAllMeetingNicknameByCCServerRet queryAllMeetingNicknameByCCServer()

	//记录用户的会议记录
	//userId		用户ID
	//callId		会议ID
	//enterTime		进入会议时间( 毫秒级时间戳 ）
	//quitTime		退出会议时间（毫秒级时间戳）
	//flowerNum		收到的花朵数量 ，不为空，没有的话传0
	common.ComRet addMeetingCallRecordByCCServer(1: i64 userId, 2: string callId, 3: i64 enterTime, 4: i64 quitTime, 5: i32 flowerNum)

	//批量获取用户设备信息
	//userIds	用户ID, 以逗号做分割符
	QueryUserDeviceInfoByCCServerRet queryUserDeviceInfoByCCServer(1: string userIds)	

	//根据话题类型获取话题列表
	//topicType		话题类型 
	QueryTopicListByCCServerRet queryTopicListByCCServer(1: i32 topicType)

	//获取超能力标签配置列表
	QueryAllSuperPowerTagListByCCServerRet queryAllSuperPowerTagListByCCServer()
	
	//是否可以踢人
	//userId	用户ID
	IsCanKickOutRet isCanKickOut(1: i64 userId)
	
	//批量获取用户信息
	//userId	用户ID
	//friendUserIds	其他用户ID, 多个用户以逗号做分割符
	QueryUserListByMeetingRet queryUserListByMeeting(1: i64 userId, 2: string friendUserIds)
	
	//判断用户是否在随机呼白名单内
	//userId	用户ID
	//return    0: 不在　　1=在
	IsUserInRandomCallingWhiteListRet isUserInRandomCallingWhiteList(1: i64 userId)
    //判断用户是否外挂用户
    SpamUserLevel queryUserSpamLevel(1: i64 userId);
}


