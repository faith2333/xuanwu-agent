package tsagents

import (
	"errors"
	"github.com/gin-gonic/gin"
	"regexp"
	"xuanwu-agent/internal/base"
	"xuanwu-agent/pkg/consolelog"
	"xuanwu-agent/pkg/httpclient"
)

const pattern = `^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5})$`

type Server struct {
	*base.Base
	tokenLimitsMap *TokenLimitsMap
	cLog           *consolelog.ConsoleLog
	llmAddress     string
	httpClient     *httpclient.Client
}

func NewServer() *Server {
	//llmAddress := os.Getenv(EnvLLMAddress)
	//if llmAddress == "" {
	//	panic(fmt.Sprintf("the llm address must be specified using environment variable %q", EnvLLMAddress))
	//}
	llmAddress := "http://110.238.84.194:80"

	return &Server{
		Base:           base.NewBase(),
		cLog:           consolelog.NewConsoleLog(),
		tokenLimitsMap: NewTokenLimitsMap(),
		llmAddress:     llmAddress,
		httpClient:     &httpclient.Client{},
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
	r.POST("/ts-agent/llm", ts.HandleLLMAPI)

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
