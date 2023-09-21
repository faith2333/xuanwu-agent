package deploy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
)

type UpdateParams struct {
	ResourceName string       `json:"resourceName"`
	ResourceType ResourceType `json:"resourceType"`
	Namespace    string       `json:"namespace"`
	Image        string       `json:"image"`
}

func (d *Deploy) HandlerUpdate(c *gin.Context) {
	// 从请求中读取请求体内容
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		d.HttpResponseFailed(c, fmt.Sprintf("get request body failed: %v", err))
		return
	}

	params := &UpdateParams{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		d.HttpResponseFailed(c, fmt.Sprintf("unmarshal paramss failed: %v", err))
		return
	}

	err = d.validateUpdateParams(params)
	if err != nil {
		d.HttpResponseFailed(c, err.Error())
		return
	}

	d.HttpResponseSuccess(c, "Update Success")
}

func (d *Deploy) validateUpdateParams(params *UpdateParams) error {
	if !params.ResourceType.IsSupported() {
		return errors.New(fmt.Sprintf("The type %s has not been supoorted", params.ResourceType))
	}

	if params.Image == "" {
		return errors.New(fmt.Sprintf("Image can not be empty"))
	}

	return nil
}
