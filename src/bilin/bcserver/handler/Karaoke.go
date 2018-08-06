package handler

import (
	"bilin/bcserver/bccommon"
	"bilin/bcserver/domain/adapter"
	"bilin/bcserver/domain/collector"
	"bilin/bcserver/domain/entity"
	"bilin/bcserver/domain/service"
	"bilin/protocol"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	SONGLISTLIMIT = 20 // 最多只能点20首歌

	NORMARL_TERMINATE = 0 // 正常结束
)

// 主持人开启/关闭K歌功能
func (this *BCServantObj) KaraokeOperation(ctx context.Context, req *bilin.KaraokeOperationReq) (resp *bilin.KaraokeOperationResp, err error) {
	const prefix = "KaraokeOperation "
	roomid := req.Header.Roomid
	userid := req.Header.Userid
	resp = &bilin.KaraokeOperationResp{Commonret: bccommon.SUCCESSMESSAGE}
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomid, userid); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp))
		return resp, nil
	}

	//检查用户权限
	if user.Role != entity.ROLE_HOST {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_OPT_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error(prefix+"KaraokeOperation failed, permission denied", zap.Any("User", userid), zap.Any("Role", user.Role))
		return
	}

	if room.Karaokeswitch == req.Opt {
		log.Info("[-]KaraokeOperation room.Karaokeswitch not change", zap.Any("req", req))
		return
	}

	if req.Opt == bilin.BaseRoomInfo_CLOSEKARAOKE {
		collector.CloseKaraoke(room)
	} else {
		collector.OpenKaraoke(room)
	}

	//更新房间信息
	collector.StorageRoomInfo(room)

	//推送广播消息
	adapter.PushBaseRoomInfoToRoom(room)

	log.Info("[+]"+prefix+"success", zap.Any("req", req), zap.Any("resp", resp))
	return
}

// 添加歌曲，只有嘉宾和主持人能操作
func (this *BCServantObj) KaraokeAddSong(ctx context.Context, req *bilin.KaraokeAddSongReq) (resp *bilin.KaraokeAddSongResp, err error) {
	const prefix = "KaraokeAddSong "
	roomid := req.Header.Roomid
	userid := req.Header.Userid
	resp = &bilin.KaraokeAddSongResp{Commonret: bccommon.SUCCESSMESSAGE}
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomid, userid); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp))
		return resp, nil
	}

	if room.Karaokeswitch != bilin.BaseRoomInfo_OPENKARAOKE {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SWTICH_CLOSED, fmt.Sprintf("请先打开K歌开关"))
		log.Error(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp))
		return resp, nil
	}

	//检查用户权限,只有麦上的用户才能点歌
	if ifOnMike, _ := service.RedisIfUserOnMike(roomid, userid); !ifOnMike {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_OPT_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error(prefix+" failed", zap.Any("resp", resp))
		return resp, nil
	}

	if collector.GetSongsCount(roomid) >= SONGLISTLIMIT {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_ADD_FAILED, fmt.Sprintf("房间已点歌曲超过20首，无法点歌，请稍后再试"))
		log.Error(prefix+" failed", zap.Any("resp", resp))
		return resp, nil
	}

	song := &bilin.KaraokeSongInfo{Id: collector.GenXid(), SongName: req.SongName, Resourceid: req.Resourceid, Userinfo: collector.LocalUserToSendInfo(user), Status: bilin.KaraokeSongInfo_PREPARE}
	coErr := collector.AddSong(room, song)
	if coErr != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_ADD_FAILED, fmt.Sprintf("添加歌曲失败"))
		log.Error(prefix+" failed", zap.Any("resp", resp))
		return resp, nil
	}

	//推送广播消息
	adapter.PushAddSongToRoom(room, song)
	adapter.PushSongsListToRoom(room)

	log.Info("[+]"+prefix+"success", zap.Any("req", req), zap.Any("resp", resp), zap.Any("song", song))
	return
}

// 主持人开始播放歌曲， 嘉宾也可以从暂停状态开始歌曲
func (this *BCServantObj) KaraokeStartSing(ctx context.Context, req *bilin.KaraokeStartSingReq) (resp *bilin.KaraokeStartSingResp, err error) {
	const prefix = "KaraokeStartSing "
	roomid := req.Header.Roomid
	userid := req.Header.Userid
	resp = &bilin.KaraokeStartSingResp{Commonret: bccommon.SUCCESSMESSAGE}
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomid, userid); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error(prefix+" failed", zap.Any("resp", resp))
		return resp, nil
	}

	if room.Karaokeswitch != bilin.BaseRoomInfo_OPENKARAOKE {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SWTICH_CLOSED, fmt.Sprintf("请先打开K歌开关"))
		log.Error(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp))
		return resp, nil
	}

	displaySong := collector.GetDisplaySong(room)
	if displaySong == nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_NOT_EXIST, fmt.Sprintf("歌曲不存在"))
		log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp))
		return
	}

	if displaySong.Status == bilin.KaraokeSongInfo_SINGING { //只有暂停和准备状态可以开始
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_ALREADY_IN_SINGING, fmt.Sprintf("歌曲已经在播放中"))
		log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp))
		return
	}

	//检查用户权限
	if user.Role != entity.ROLE_HOST {
		if user.UserID != displaySong.Userinfo.Userid { //不是主持人，也不是嘉宾
			resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_OPT_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
			log.Error(prefix, zap.Any("User", userid), zap.Any("Role", user.Role))
			return
		} else { //嘉宾点的歌曲
			if displaySong.Status != bilin.KaraokeSongInfo_PAUSE { // 是嘉宾，但是歌曲状态不是暂停，也不能操作
				resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_OPT_NO_RIGHT, fmt.Sprintf("用户没有权限"))
				log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp))
				return
			}
		}
	}

	//更新歌曲状态
	displaySong.Status = bilin.KaraokeSongInfo_SINGING
	coErr := collector.ChangeDisplaySongStatus(room, displaySong)
	if coErr != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_CHANGESTATUS_FAILED, fmt.Sprintf("改变歌曲状态失败"))
		log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp))
		return
	}

	//推送广播消息
	adapter.PushStartSingToRoom(room, displaySong)

	log.Info("[+]"+prefix+"success", zap.Any("req", req), zap.Any("resp", resp), zap.Any("song", displaySong))
	return
}

// 主持人对某个歌曲置顶
func (this *BCServantObj) KaraokeSongSetTop(ctx context.Context, req *bilin.KaraokeSongSetTopReq) (resp *bilin.KaraokeSongSetTopResp, err error) {
	const prefix = "KaraokeSongSetTop "
	roomid := req.Header.Roomid
	userid := req.Header.Userid
	resp = &bilin.KaraokeSongSetTopResp{Commonret: bccommon.SUCCESSMESSAGE}
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomid, userid); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error(prefix+" failed", zap.Any("resp", resp))
		return resp, nil
	}

	if room.Karaokeswitch != bilin.BaseRoomInfo_OPENKARAOKE {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SWTICH_CLOSED, fmt.Sprintf("请先打开K歌开关"))
		log.Error(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp))
		return resp, nil
	}

	//检查用户权限
	if user.Role != entity.ROLE_HOST {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_OPT_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error(prefix, zap.Any("User", userid), zap.Any("Role", user.Role))
		return
	}

	song := collector.GetSongInfoByID(room, req.Songid)
	if song == nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_NOT_EXIST, fmt.Sprintf("歌曲不存在"))
		log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp))
		return
	}
	coErr := collector.SetTopSong(room, song)
	if coErr != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_SETTOP_FAILED, fmt.Sprintf("歌曲置顶错误"))
		log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp))
		return
	}

	//给被置顶的用户发送单播通知
	adapter.PushNotifyToUser(room.Roomid, []int64{int64(song.Userinfo.Userid)}, &bilin.SongSetTopNotify{Song: song}, bilin.MinType_BC_NotifySongSetTop)

	//给频道内的用户发送歌曲列表的广播
	adapter.PushSongsListToRoom(room)

	log.Info("[+]"+prefix+"success", zap.Any("req", req), zap.Any("resp", resp), zap.Any("song", song))
	return
}

// 主持人或者嘉宾删除歌曲
func (this *BCServantObj) KaraokeDelSong(ctx context.Context, req *bilin.KaraokeDelSongReq) (resp *bilin.KaraokeDelSongResp, err error) {
	const prefix = "KaraokeDelSong "
	roomid := req.Header.Roomid
	userid := req.Header.Userid
	resp = &bilin.KaraokeDelSongResp{Commonret: bccommon.SUCCESSMESSAGE}
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomid, userid); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error(prefix+" failed", zap.Any("resp", resp))
		return resp, nil
	}

	if room.Karaokeswitch != bilin.BaseRoomInfo_OPENKARAOKE {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SWTICH_CLOSED, fmt.Sprintf("请先打开K歌开关"))
		log.Error(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp))
		return resp, nil
	}

	song := collector.GetSongInfoByID(room, req.Songid)
	if song == nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_NOT_EXIST, fmt.Sprintf("歌曲不存在"))
		log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp))
		return
	}
	if user.Role != entity.ROLE_HOST && user.UserID != song.Userinfo.Userid { //主持人可以删除所有的歌曲，嘉宾只能删除自己点的歌曲
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_OPT_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp))
		return resp, nil
	}

	//正在播放或者暂停的歌曲不能删除
	if song.Status != bilin.KaraokeSongInfo_PREPARE {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_DEL_FAILED, fmt.Sprintf("正在播放的歌曲不能删除"))
		log.Error(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp), zap.Any("song", song))
		return resp, nil
	}

	if coErr := collector.DelSong(room, song); coErr != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_DEL_FAILED, fmt.Sprintf("删除歌曲失败"))
		log.Error(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp), zap.Any("song", song))
		return resp, nil
	}

	//被别人删掉才需要发push
	if song.Userinfo.Userid != userid {
		adapter.PushNotifyToUser(room.Roomid, []int64{int64(song.Userinfo.Userid)}, &bilin.DelSongNotify{Song: song}, bilin.MinType_BC_NotifyDelSong)
	}

	//给频道内的用户发送歌曲列表的广播
	adapter.PushSongsListToRoom(room)

	log.Info("[+]"+prefix+"success", zap.Any("req", req), zap.Any("resp", resp), zap.Any("song", song))
	return
}

// 主持人或者嘉宾暂停歌曲
func (this *BCServantObj) KaraokePauseSong(ctx context.Context, req *bilin.KaraokePauseSongReq) (resp *bilin.KaraokePauseSongResp, err error) {
	const prefix = "KaraokePauseSong "
	roomid := req.Header.Roomid
	userid := req.Header.Userid
	resp = &bilin.KaraokePauseSongResp{Commonret: bccommon.SUCCESSMESSAGE}
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomid, userid); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error(prefix+" failed", zap.Any("resp", resp))
		return resp, nil
	}

	if room.Karaokeswitch != bilin.BaseRoomInfo_OPENKARAOKE {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SWTICH_CLOSED, fmt.Sprintf("请先打开K歌开关"))
		log.Error(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp))
		return resp, nil
	}

	//主持人或者嘉宾暂停歌曲
	song := collector.GetDisplaySong(room)
	if song == nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_NOT_EXIST, fmt.Sprintf("歌曲不存在"))
		log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp))
		return
	}
	if user.Role != entity.ROLE_HOST && user.UserID != song.Userinfo.Userid {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_OPT_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error(prefix+" failed", zap.Any("resp", resp), zap.Any("song", song))
		return resp, nil
	}
	if song.Status != bilin.KaraokeSongInfo_SINGING {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_NOT_START, fmt.Sprintf("只有在播放状态才能暂停"))
		log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp), zap.Any("song", song))
		return
	}

	song.Status = bilin.KaraokeSongInfo_PAUSE
	//更新歌曲状态
	coErr := collector.ChangeDisplaySongStatus(room, song)
	if coErr != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_CHANGESTATUS_FAILED, fmt.Sprintf("改变歌曲状态失败"))
		log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp))
		return
	}

	//给频道内的用户发送歌曲列表的广播
	adapter.PushPauseSongToRoom(room, song)

	log.Info("[+]"+prefix+"success", zap.Any("req", req), zap.Any("resp", resp), zap.Any("song", song))
	return
}

// 主持人或者嘉宾结束歌曲
func (this *BCServantObj) KaraokeTerminateSong(ctx context.Context, req *bilin.KaraokeTerminateSongReq) (resp *bilin.KaraokeTerminateSongResp, err error) {
	const prefix = "KaraokeTerminateSong "
	roomid := req.Header.Roomid
	userid := req.Header.Userid
	resp = &bilin.KaraokeTerminateSongResp{Commonret: bccommon.SUCCESSMESSAGE}
	log.Info(prefix+"begin", zap.Any("req", req))

	defer func(now time.Time) {
		httpmetrics.DefReport(strings.TrimSpace(prefix), int64(resp.Commonret.Ret), now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	var room *entity.Room
	var user *entity.User
	if room, user, err = this.CommonCheckAuth(roomid, userid); err != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_ILLEGAL_MESSAGE, bccommon.COMMERRORILLEGALDESC)
		log.Error(prefix+" failed", zap.Any("resp", resp))
		return resp, nil
	}

	if room.Karaokeswitch != bilin.BaseRoomInfo_OPENKARAOKE {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SWTICH_CLOSED, fmt.Sprintf("请先打开K歌开关"))
		log.Error(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp))
		return resp, nil
	}

	//主持人或者嘉宾结束歌曲
	song := collector.GetDisplaySong(room)
	if song == nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_NOT_EXIST, fmt.Sprintf("歌曲不存在"))
		log.Error(prefix, zap.Any("req", req), zap.Any("resp", resp))
		return
	}
	if user.Role != entity.ROLE_HOST && user.UserID != song.Userinfo.Userid {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_OPT_NO_RIGHT, fmt.Sprintf("你没有该权限哦~"))
		log.Error(prefix+" failed", zap.Any("resp", resp))
		return resp, nil
	}

	//准备状态下的歌曲无法结束，只能删除
	if song.Status == bilin.KaraokeSongInfo_PREPARE {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_TERMINATE_FAILED, fmt.Sprintf("准备中的歌曲无法结束"))
		log.Error(prefix+" failed", zap.Any("req", req), zap.Any("resp", resp), zap.Any("song", song))
		return resp, nil
	}

	coErr := collector.DelSong(room, song)
	if coErr != nil {
		resp.Commonret = bccommon.UserDefinedFailed(bilin.CommonRetInfo_KARAOKE_SONG_TERMINATE_FAILED, fmt.Sprintf("服务器开小差了，再试试呗~"))
		log.Error(prefix+" failed", zap.Any("resp", resp), zap.Any("song", song))
		return resp, nil
	}

	//广播通知,异常情况下才会通知
	if req.Flag != NORMARL_TERMINATE {
		adapter.PushTerminateSongToRoom(room, userid, song)
	}

	//给频道内的用户发送歌曲列表的广播
	adapter.PushSongsListToRoom(room)

	log.Info("[+]"+prefix+"success", zap.Any("req", req), zap.Any("resp", resp), zap.Any("song", song))
	return
}
