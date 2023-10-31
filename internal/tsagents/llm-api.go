package tsagents

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"xuanwu-agent/pkg/httpclient"
)

const (
	EnvLLMAddress = "ENV_LLM_ADDRESS"
)

type LLMAPIReqBody struct {
	WorkflowID string      `json:"WorkflowID"`
	URL        string      `json:"URL"`
	Data       interface{} `json:"Data"`
}

func (ts *Server) HandleLLMAPI(c *gin.Context) {
	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ts.HttpResponseFailed(c, fmt.Sprintf("get request body failed: %v", err))
		return
	}

	params := &LLMAPIReqBody{}
	err = json.Unmarshal(reqBody, &params)
	if err != nil {
		ts.HttpResponseFailed(c, fmt.Sprintf("unmarshal request body failed: %v", err))
		return
	}

	if ok := ts.tokenLimitsMap.Add(c.ClientIP()); !ok {
		ts.HttpResponseFailed(c, fmt.Sprintf("IP: %q has been blocked due to exhaustion token", c.ClientIP()))
		return
	}

	resp, err := ts.dialLLM(params.URL, params.Data)
	if err != nil {
		ts.HttpResponseFailed(c, fmt.Sprintf("dial LLM failed: %v", err))
		return
	}
	var data = make(map[string]interface{})
	err = json.Unmarshal(resp, &data)
	if err != nil {
		ts.HttpResponseFailed(c, fmt.Sprintf("get response data failed: %v", err))
		return
	}

	ts.HttpResponseSuccess(c, data)
	go func() {
		err = ts.InsertRecord(params.WorkflowID, "recruitment", data)
		if err != nil {
			ts.cLog.Errorf("insert execution record failed %v", err)
		}
	}()
	return
}

func (ts *Server) dialLLM(url string, data interface{}) ([]byte, error) {
	return ts.httpClient.WithMethod(httpclient.MethodPOST).
		WithURL(ts.llmAddress + url).WithBody(data).WithContentTypeJSON().Do()
}
