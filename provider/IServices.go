package service

import (
	"github.com/taylor840326/aicloud-sdk-go/ackcs"
	"github.com/taylor840326/aicloud-sdk-go/common/errors"
)

type IServices interface {
	CreateService(request ackcs.CreateServicesRequest) errors.AIErrors
	DeleteService(serviceUuid string) errors.AIErrors
	DescribeServices(serviceUuid string) (ackcs.DescribeServicesResponse, errors.AIErrors)
	DescribeServiceJupyter(serviceUuid string) (ackcs.DescribeServicesJupyterResponse, errors.AIErrors)
	DescribeServiceSSH(serviceUuid string) (ackcs.DescribeServicesSSHResponse, errors.AIErrors)
}
