package tsagents

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"sync"
	"time"
	"xuanwu-agent/internal/base"
)

const (
	EnvDBInfo = "ENV_DB_INFO"
)

type TSAgentExecutionRecords struct {
	ID         int64        `json:"id" gorm:"primaryKey;autoIncrement"`
	Type       string       `json:"type" gorm:"type:varchar(64)"`
	WorkflowID string       `json:"workflowID" gorm:"type:varchar(64)"`
	Data       base.TypeMap `json:"data" gorm:"type:json"`
	CreateUser string       `json:"createUser" gorm:"type:varchar(16)"`
	ModifyUser string       `json:"modifyUser" gorm:"type:varchar(16)"`
	GmtCreate  string       `json:"gmtCreate" gorm:"type:varchar(64)"`
	GmtModify  string       `json:"gmtModify" gorm:"type:varchar(64)"`
	Deleted    int64        `json:"deleted" gorm:"type:int(2)"`
}

func (records *TSAgentExecutionRecords) TableName() string {
	return "ta_execution_records"
}

var runOnce = &sync.Once{}

func newDB() (*gorm.DB, error) {
	dbInfo := os.Getenv(EnvDBInfo)
	if dbInfo == "" {
		return nil, errors.New("get database information from environment failed: not found")
	}

	db, err := gorm.Open(mysql.Open(dbInfo))
	if err != nil {
		return nil, err
	}
	runOnce.Do(func() {
		err = db.AutoMigrate(TSAgentExecutionRecords{})
		if err != nil {
			panic(err)
		}
	})

	return db, err
}

func (ts *Server) InsertRecord(workflowID, workflowType string, data map[string]interface{}) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	record := &TSAgentExecutionRecords{
		Type:       workflowType,
		WorkflowID: workflowID,
		CreateUser: "tsagent",
		ModifyUser: "tsagent",
		Data:       base.NewTypeMap(data),
		GmtCreate:  now,
		GmtModify:  now,
		Deleted:    0,
	}

	err := ts.db.Model(&TSAgentExecutionRecords{}).Create(&record).Error
	if err != nil {
		return err
	}

	return nil
}

type PageInfo struct {
	PageIndex int   `json:"pageIndex"`
	PageSize  int   `json:"pageSize"`
	Total     int64 `json:"total"`
}

func (ts *Server) ListRecords(pageIndex, pageSize int, workflowType, workflowID string) (*PageInfo, []*TSAgentExecutionRecords, error) {
	offSet := (pageIndex - 1) * pageSize
	records := make([]*TSAgentExecutionRecords, 0)

	query := ts.db.Model(&TSAgentExecutionRecords{}).Where("type = ? and workflow_id = ? and deleted = 0", workflowType, workflowID)

	var total int64 = 0
	err := query.Count(&total).Error
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("execute db query failed: %v", err))
	}
	err = query.Offset(offSet).Limit(pageSize).Find(&records).Error
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("execute db query failed: %v", err))
	}

	return &PageInfo{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     total,
	}, records, nil
}

func (ts *Server) DeleteRecords(ids []int64) error {
	db := ts.db.Model(&TSAgentExecutionRecords{})

	for _, id := range ids {
		err := db.Where("id = ?", id).Update("deleted", id).Error
		if err != nil {
			return err
		}
	}
	return nil
}
