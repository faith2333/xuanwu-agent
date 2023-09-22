package deploy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	err = d.update(params)
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

	if params.Namespace == "" {
		return errors.New("namespace can not be empty")
	}

	if params.ResourceName == "" {
		return errors.New("resource name can not be empty")
	}

	return nil
}

func (d *Deploy) update(param *UpdateParams) error {
	switch param.ResourceType.Upper() {
	case ResourceTypeDeployment:
		return d.updateDeployment(param)
	}
	return errors.New(fmt.Sprintf("resource type %s has not been supported", param.ResourceType))
}

func (d *Deploy) updateDeployment(params *UpdateParams) error {
	ctx := context.Background()
	deployment, err := d.kubeClient.AppsV1().Deployments(params.Namespace).Get(ctx, params.ResourceName, metav1.GetOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("get the deployment %q in namespace %q failed: %v", params.ResourceName, params.Namespace, err))
	}

	deployment.Spec.Template.Spec.Containers[0].Image = params.Image

	_, err = d.kubeClient.AppsV1().Deployments(params.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
