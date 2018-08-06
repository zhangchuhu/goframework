package updater

import (
	"bilin/searchserver/config"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/Shopify/sarama"
)

func TestUserToURL(t *testing.T) {
	var (
		err  error
		l    *url.URL
		rsp  *http.Response
		body []byte
	)

	user := UserU{
		Id:      "40373825",
		BilinId: "52851301",
		Name:    "小伙伴",
	}

	if l, err = user.ToURL("bilin_user_update_test"); err != nil {
		t.Fatalf("convert to url fail: %v", err)
	}
	t.Logf("url: %v", l)

	if rsp, err = http.Get(l.String()); err != nil {
		t.Fatalf("http update fail: %v", err)
	}
	defer rsp.Body.Close()

	if body, err = ioutil.ReadAll(rsp.Body); err != nil {
		t.Fatalf("read http response fail: %v", err)
	}
	t.Logf("http update got: %v", string(body))
}

func TestRoomToURL(t *testing.T) {
	var (
		err  error
		l    *url.URL
		rsp  *http.Response
		body []byte
	)

	room := RoomU{
		Id:        "410316952",
		DisplayId: "10086",
		Name:      "大合唱",
		Live:      "1",
	}

	if l, err = room.ToURL("bilin_room_update_test"); err != nil {
		t.Fatalf("convert to url fail: %v", err)
	}
	t.Logf("url: %v", l)

	if rsp, err = http.Get(l.String()); err != nil {
		t.Fatalf("http update fail: %v", err)
	}
	defer rsp.Body.Close()

	if body, err = ioutil.ReadAll(rsp.Body); err != nil {
		t.Fatalf("read http response fail: %v", err)
	}
	t.Logf("http update got: %v", string(body))
}

func TestSongToURL(t *testing.T) {
	var (
		err  error
		l    *url.URL
		rsp  *http.Response
		body []byte
	)

	song := SongU{
		Id:     "1",
		Name:   "绅士(没打包)",
		Artist: "薛之谦",
	}

	if l, err = song.ToURL("bilin_song_update_test"); err != nil {
		t.Fatalf("convert to url fail: %v", err)
	}
	t.Logf("url: %v", l)

	if rsp, err = http.Get(l.String()); err != nil {
		t.Fatalf("http update fail: %v", err)
	}
	defer rsp.Body.Close()

	if body, err = ioutil.ReadAll(rsp.Body); err != nil {
		t.Fatalf("read http response fail: %v", err)
	}
	t.Logf("http update got: %v", string(body))
}

func TestKafkaProduce(t *testing.T) {
	c := sarama.NewConfig()
	c.Producer.Return.Successes = true
	p, err := sarama.NewSyncProducer(config.GetAppConfig().KafkaBrokers, c)
	if err != nil {
		t.Fatalf("can not create kafka producer: %v", err)
	}

	user := UserU{
		Id:      "40373825",
		BilinId: "52851301",
		Name:    "小伙伴",
	}
	bytes, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("can not marshal: %v", err)
	}

	_, _, err = p.SendMessage(&sarama.ProducerMessage{
		Topic: "bilin_user_update_test",
		Value: sarama.ByteEncoder(bytes),
	})
	if err != nil {
		t.Fatalf("can not send kafka message: %v", err)
	}
}
