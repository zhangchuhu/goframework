package index_test

import (
	"bilin/apigateway/common/index"
	"testing"
)

func TestIndex(t *testing.T) {
	var i index.Index
	i.InitCache()
	i.Add(999, "", "")
	ok := i.Match(999, "1.0.0")
	t.Logf("999 match:%v", ok)
	ok = i.Match(888, "1.0.0")
	t.Logf("888 match:%v", ok)
	i.Add(777, "1.0.0,2.0.0,3.0.0", "")
	ok = i.Match(777, "1.0.0")
	t.Logf("777 match:%v", ok)
	ok = i.Match(777, "2.0.0")
	t.Logf("777 match:%v", ok)
	ok = i.Match(777, "3.0.0")
	t.Logf("777 match:%v", ok)
	ok = i.Match(777, "4.0.0")
	t.Logf("777 match:%v", ok)
	ok = i.Match(666, "1.0.0")
	t.Logf("666 match:%v", ok)
	ok = i.Match(666, "4.0.0")
	t.Logf("666 match:%v", ok)
}
