package tsagents

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
)

type QueryParams struct {
	PageIndex  int    `json:"pageIndex"`
	PageSize   int    `json:"pageSize"`
	WorkflowID string `json:"workflowID"`
}

func (ts *Server) ServerRecords(c *gin.Context) {
	params, err := getListRecordsQueryParams(c)
	if err != nil {
		ts.HttpResponseFailed(c, err.Error())
		return
	}

	pageInfo, records, err := ts.ListRecords(params.PageIndex, params.PageSize, "recruitment", params.WorkflowID)
	if err != nil {
		ts.HttpResponseFailed(c, err.Error())
		return
	}

	ts.HttpResponseSuccess(c, struct {
		PageInfo *PageInfo                  `json:"pageInfo"`
		Records  []*TSAgentExecutionRecords `json:"records"`
	}{
		PageInfo: pageInfo,
		Records:  records,
	})
}

type DeleteRecordsParams struct {
	IDs []int64 `json:"ids"`
}

func (ts *Server) ServerDeleteRecords(c *gin.Context) {
	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ts.HttpResponseFailed(c, fmt.Sprintf("get request body failed: %v", err))
		return
	}

	params := &DeleteRecordsParams{}
	err = json.Unmarshal(reqBody, &params)
	if err != nil {
		ts.HttpResponseFailed(c, fmt.Sprintf("unmarshal request body failed: %v", err))
		return
	}

	err = ts.DeleteRecords(params.IDs)
	if err != nil {
		ts.HttpResponseFailed(c, err.Error())
		return
	}

	ts.HttpResponseSuccess(c, "")
}

func getListRecordsQueryParams(c *gin.Context) (*QueryParams, error) {
	params := &QueryParams{}
	var err error

	stringPageIndex := c.DefaultQuery("pageIndex", "1")
	params.PageIndex, err = strconv.Atoi(stringPageIndex)
	if err != nil {
		return nil, err
	}

	stringPageSize := c.DefaultQuery("pageSize", "20")
	params.PageSize, err = strconv.Atoi(stringPageSize)
	if err != nil {
		return nil, err
	}

	params.WorkflowID = c.Query("workflowID")
	return params, nil
}
