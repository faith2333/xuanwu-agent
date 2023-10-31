package tsagents

import (
	"github.com/gin-gonic/gin"
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
