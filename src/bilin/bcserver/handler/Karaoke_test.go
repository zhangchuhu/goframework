package handler

import (
	"bilin/protocol"
	"context"
	"testing"
)

func TestKaraokeOpenKaraoke(t *testing.T) {
	resp, err := s.KaraokeOperation(context.TODO(), &bilin.KaraokeOperationReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Opt: bilin.BaseRoomInfo_OPENKARAOKE,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestKaraokeCloseKaraoke(t *testing.T) {
	resp, err := s.KaraokeOperation(context.TODO(), &bilin.KaraokeOperationReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Opt: bilin.BaseRoomInfo_CLOSEKARAOKE,
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestKaraokeHostAddSong(t *testing.T) {
	resp, err := s.KaraokeAddSong(context.TODO(), &bilin.KaraokeAddSongReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		SongName:   "host test",
		Resourceid: "aaaaa",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestKaraokeGuestAddSong(t *testing.T) {
	resp, err := s.KaraokeAddSong(context.TODO(), &bilin.KaraokeAddSongReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: uid3,
		},
		SongName:   "uid3 test",
		Resourceid: "ddddd",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestKaraokeStartSing(t *testing.T) {
	resp, err := s.KaraokeStartSing(context.TODO(), &bilin.KaraokeStartSingReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Songid: "bcplikgmm5fo5aui87q0",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestKaraokeSongSetTop(t *testing.T) {
	resp, err := s.KaraokeSongSetTop(context.TODO(), &bilin.KaraokeSongSetTopReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Songid: "bcplikgmm5fo5aui87q0",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestKaraokePauseSong(t *testing.T) {
	resp, err := s.KaraokePauseSong(context.TODO(), &bilin.KaraokePauseSongReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Songid: "bcplo7omm5foa8h4tcs0",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestKaraokeDelSong(t *testing.T) {
	resp, err := s.KaraokeDelSong(context.TODO(), &bilin.KaraokeDelSongReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Songid: "bcplcb0mm5fvh5g7b1l0",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}

func TestKaraokeTerminateSong(t *testing.T) {
	resp, err := s.KaraokeTerminateSong(context.TODO(), &bilin.KaraokeTerminateSongReq{
		Header: &bilin.Header{
			Roomid: roomid,
			Userid: hostuid,
		},
		Songid: "1111",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", *resp)

}
