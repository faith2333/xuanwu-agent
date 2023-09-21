package deploy

import (
	"github.com/gin-gonic/gin"
	"xuanwu-agent/internal/base"
	"xuanwu-agent/pkg/consolelog"
)

type Deploy struct {
	*base.Base
	cLog *consolelog.ConsoleLog
}

func NewDeploy() *Deploy {
	return &Deploy{
		Base: base.NewBase(),
		cLog: consolelog.NewConsoleLog(),
	}
}

func (d *Deploy) Listen(addr string) int {
	r := gin.Default()
	r.GET("/deploy/update", d.HandlerUpdate)

	for {
		err := r.Run(addr)
		if err != nil {
			d.cLog.Error(err.Error())
		}
	}
}
