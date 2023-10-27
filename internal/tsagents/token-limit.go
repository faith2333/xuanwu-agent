package tsagents

import (
	"os"
	"strconv"
	"sync"
	"time"
	"xuanwu-agent/pkg/consolelog"
)

const (
	DefaultTokenLimit = 3000
	ENVTokenLimit     = "TOKEN_LIMIT"
)

type TokenLimitsMap struct {
	lock   *sync.RWMutex
	Map    map[string]*TokenDetail `json:"map"`
	Limits int                     `json:"limits"`
}

type TokenDetail struct {
	count      int
	createTime time.Time
}

func NewTokenLimitsMap() *TokenLimitsMap {
	limitsEnv := os.Getenv(ENVTokenLimit)
	limits, err := strconv.Atoi(limitsEnv)
	if err != nil {
		(&consolelog.ConsoleLog{}).Waringf("Get token limit from env %q failed: %v, use default 3000", ENVTokenLimit, err)
		limits = DefaultTokenLimit
	}

	return &TokenLimitsMap{
		lock:   &sync.RWMutex{},
		Map:    make(map[string]*TokenDetail),
		Limits: limits,
	}
}

func (token *TokenLimitsMap) Add(clientIp string) bool {
	token.lock.Lock()
	defer token.lock.Unlock()

	detail, ok := token.Map[clientIp]
	if !ok {
		token.Map[clientIp] = &TokenDetail{
			count:      1,
			createTime: time.Now(),
		}
		return true
	}

	if duration := time.Now().Sub(detail.createTime); duration.Hours() > 24 {
		token.Map[clientIp] = &TokenDetail{
			count:      1,
			createTime: time.Now(),
		}
		return true
	}

	if detail.count >= token.Limits {
		return false
	}

	detail.count++
	return true
}
