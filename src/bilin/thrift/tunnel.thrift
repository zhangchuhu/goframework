namespace cpp com.bilin.tunnel.service

// appid = 1 营收服务appid
enum AppidType {
	REVENUE_SRV = 1,
	NEW_BCSERVER = 11,
}

service Tunnel {
  /***
     功能：对指定用户发送消息
     参数：
        appid：业务调用方唯一标志，统一分配
        uid：单播用户的唯一id
        msg：业务方的透传消息
  ***/
  i32 unicastByUid(1: i64 appid, 2: i32 uid, 3: string msg);

  /***
     功能：对指定用户发送消息
     参数：
        appid：业务调用方唯一标志，统一分配
        uid：单播用户的唯一id
        msg：业务方的透传消息
  ***/
  i32 unicastByUidEx(1: i64 appid, 2: i32 uid, 3: string msg 4: i32 msg_type);

/***
   功能：对指定房间的用户发送消息
   参数：
      appid：业务调用方唯一标志，统一分配
      sid: 频道id
      uid：单播用户的唯一id
      msg：业务方的透传消息
***/
i32 unicastToRoomByUidEx(1: i64 appid, 2: i32 sid, 3: i32 uid, 4: string msg 5: i32 msg_type);

  /****
     功能：对指定频道发送广播消息
     参数：
        appid：业务调用方唯一标志，统一分配
        sid：广播频道id
        msg：业务方的透传消息
        msg_type:业务对应的消息类型
  ****/
  i32 broadcastBySid(1: i64 appid,2: i32 sid, 3: string msg);

  /****
     功能：对指定频道发送广播消息
     参数：
        appid：业务调用方唯一标志，统一分配
        sid：广播频道id
        msg：业务方的透传消息
  ****/
  i32 broadcastBySidEx(1: i64 appid,2: i32 sid,3: i32 msg_type,  4: string msg);

  /****
	功能：官频上麦, 自动把原来在麦上用户切下去
	参数：
	uid：要上麦主播uid
	sid：广播频道id
	msg：业务方的透传消息
	返回：
	异步回调结果
   ****/
  i32 onOfficialMike(1: i64 uid, 2: i32 sid, 3: string msg);

  /****
	功能：官频下麦, 仅仅下麦，后续可能用到
	参数：
	uid：要下麦主播uid
	sid：广播频道id
	msg：业务方的透传消息
	返回：
	异步回调结果
   ****/
  i32 offOfficialMike(1: i64 uid, 2: i32 sid, 3: string msg);

}
