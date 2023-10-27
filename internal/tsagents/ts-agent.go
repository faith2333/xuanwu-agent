package tsagents

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"regexp"
	"xuanwu-agent/internal/base"
	"xuanwu-agent/pkg/consolelog"
)

const pattern = `^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5})$`

type Server struct {
	*base.Base
	tokenLimitsMap *TokenLimitsMap
	cLog           *consolelog.ConsoleLog
	llmAddress     string
}

func NewServer() *Server {
	llmAddress := os.Getenv(EnvLLMAddress)
	if llmAddress == "" {
		panic(fmt.Sprintf("the llm address must be specified using environment variable %q", EnvLLMAddress))
	}

	return &Server{
		Base:           base.NewBase(),
		cLog:           consolelog.NewConsoleLog(),
		tokenLimitsMap: NewTokenLimitsMap(),
		llmAddress:     llmAddress,
	}
}

type FlagParams struct {
	Address string `json:"address"`
}

func (ts *Server) Listen(params *FlagParams) int {
	err := ts.init(params)
	if err != nil {
		ts.cLog.Error(err.Error())
		return base.ReturnError
	}

	r := gin.Default()
	r.POST("/ts-agent/llm")

	for {
		err = r.Run(params.Address)
		if err != nil {
			ts.cLog.Error(err.Error())
		}
	}
}

func (ts *Server) init(params *FlagParams) error {
	if ts.cLog == nil {
		ts.cLog = consolelog.NewConsoleLog()
	}

	matched, err := regexp.Match(pattern, []byte(params.Address))
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("the address is invalid, it should like 0.0.0.0:8080")
	}

	return nil
}
