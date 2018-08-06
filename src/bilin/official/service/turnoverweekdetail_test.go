package service

import (
	"context"
	"testing"
)

func TestQueryAnchorWeekPropsRecieve(t *testing.T) {
	info, err := QueryAnchorWeekPropsRecieve(context.TODO(), 6241154, "20180601000000", "20180630235959", 1, 10)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func TestQueryChannelWeekPropsRecieve(t *testing.T) {
	info, err := QueryChannelWeekPropsRecieve(context.TODO(), 17795053, "20180601000000", "20180630235959", 1, 10, 17795053)
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}
