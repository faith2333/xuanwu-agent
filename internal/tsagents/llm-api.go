package tsagents

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
	"xuanwu-agent/pkg/httpclient"
)

const (
	EnvLLMAddress = "ENV_LLM_ADDRESS"
)

type LLMAPIReqBody struct {
	WorkflowID string      `json:"WorkflowID"`
	ResumeLink string      `json:"resumeLink"`
	URL        string      `json:"url"`
	Data       interface{} `json:"data"`
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

	err = ts.initRecordData(params, data)
	if err != nil {
		ts.cLog.Error(err.Error())
	}

	ts.HttpResponseSuccess(c, data)
	return
}

func (ts *Server) dialLLM(url string, data interface{}) ([]byte, error) {
	return ts.httpClient.WithMethod(httpclient.MethodPOST).
		WithURL(ts.llmAddress + url).WithBody(data).WithContentTypeJSON().Do()
}

func (ts *Server) initRecordData(params *LLMAPIReqBody, data map[string]interface{}) error {
	if params.URL != "/check_resume_eng" {
		return nil
	}

	var candidateName string

	paramsData := make(map[string]interface{})
	dataBytes, err := json.Marshal(params.Data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dataBytes, &paramsData)
	if err != nil {
		return err
	}

	inputs, ok := paramsData["input"].([]interface{})
	if !ok {
		candidateName = "Unknown User"
		ts.cLog.Error("get input from request data failed")
	} else {
		if len(inputs) == 0 {
			candidateName = "Unknown User"
			ts.cLog.Error("get input from request data failed")
		}

		candidateName = strings.Split(inputs[0].(string), "\n")[0]
	}

	go func() {
		data["candidateName"] = candidateName
		data["resumeLink"] = params.ResumeLink
		err := ts.InsertRecord(params.WorkflowID, "recruitment", data)
		if err != nil {
			ts.cLog.Errorf("insert execution record failed %v", err)
		}
	}()

	return nil
}
