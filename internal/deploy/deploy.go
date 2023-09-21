package deploy

import (
	"errors"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"xuanwu-agent/internal/base"
	"xuanwu-agent/pkg/consolelog"
)

type Deploy struct {
	*base.Base
	cLog       *consolelog.ConsoleLog
	kubeClient kubernetes.Interface
}

func NewDeploy() *Deploy {
	return &Deploy{
		Base: base.NewBase(),
		cLog: consolelog.NewConsoleLog(),
	}
}

type FlagParams struct {
	Address   string `json:"address"`
	InCluster bool   `json:"inCluster"`
}

func (d *Deploy) Listen(flagParams *FlagParams) int {
	err := d.init(flagParams)
	if err != nil {
		d.cLog.Errorf("init deploy failed: %v", err)
		return base.ReturnError
	}

	r := gin.Default()
	r.GET("/deploy/update", d.HandlerUpdate)

	for {
		err := r.Run(flagParams.Address)
		if err != nil {
			d.cLog.Error(err.Error())
		}
	}
}

func (d *Deploy) init(flagParams *FlagParams) error {
	if d.cLog == nil {
		d.cLog = consolelog.NewConsoleLog()
	}

	if flagParams.InCluster {
		err := d.newInClusterKubeClient()
		if err != nil {
			return err
		}
	} else {
		return errors.New("currently, the deploy module of xuanwu-agent is only supported to deployed into a kubernetes cluster")
	}

	return nil
}

func (d *Deploy) newInClusterKubeClient() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	d.kubeClient = clientSet
	return nil
}
