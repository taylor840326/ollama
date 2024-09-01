package provider

import (
	"github.com/taylor840326/aicloud-sdk-go/ackcs"
)

type IServices interface {
	//类型接口
	DescribeServiceTypes() (ack_product.DescribeServiceTypesResponse, error)
	DescribePublicImages() (ack_product.DescribePublicImagesResponse, error)
	DescribeSoldOut() (ack_product.DescribeACKAvailableResourceResponse, error)
	//实例接口
	CreateService(request ackcs.CreateServicesRequest) error
	DeleteService(serviceUuid string) error
	DescribeServices(serviceUuid string) (ackcs.DescribeServicesResponse, error)
	DescribeServiceJupyter(serviceUuid string) (ackcs.DescribeServicesJupyterResponse, error)
	DescribeServiceSSH(serviceUuid string) (ackcs.DescribeServicesSSHResponse, error)
	//镜像接口
	CreateImages(request ackci.CreateImagesRequest) (ackci.CreateImagesResponse, error)
	DeleteImages(imageUuid string) error
	DescribeImages(imageUuid string) (ackci.DescribeImagesResponse, error)
}
