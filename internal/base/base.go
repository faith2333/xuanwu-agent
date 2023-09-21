package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xuanwu-agent/pkg/consolelog"
)

const (
	ReturnSuccess = 0
	ReturnError   = 1
)

type Base struct {
	cLog *consolelog.ConsoleLog
}

func NewBase() *Base {
	return &Base{
		cLog: consolelog.NewConsoleLog(),
	}
}

func (base *Base) HttpResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &gin.H{
		"code":    1000,
		"success": true,
		"message": "success",
		"data":    data,
	})
}

func (base *Base) HttpResponseFailed(c *gin.Context, mess string) {
	base.cLog.Errorf(mess)

	c.JSON(http.StatusInternalServerError, &gin.H{
		"code":    1001,
		"success": false,
		"message": mess,
		"data":    nil,
	})
}
