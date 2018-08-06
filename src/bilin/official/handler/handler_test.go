package handler

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestAuth(t *testing.T) {
	c := &gin.Context{}
	if ok, err := authCookie("mnaZDGoq8CapInm0dlvWWsDxniHqtZCeZdJfSlEQcrRZOf8oPMoxWTWGZ7aJR6DD", c); err != nil {
		t.Error(err)
	} else {
		t.Log(ok, c.GetInt64("uid"))
	}
}
