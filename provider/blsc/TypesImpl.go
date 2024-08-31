package blsc

import (
	"fmt"
	"github.com/taylor840326/aicloud-sdk-go/ack_product"
	"github.com/taylor840326/aicloud-sdk-go/common"
	"github.com/taylor840326/aicloud-sdk-go/common/errors"
)

type BlscTypesImpl struct {
	HttpProfile   common.HttpProfile   `json:"httpProfile"`
	ClientProfile common.ClientProfile `json:"clientProfile"`
}

func NewBlscTypeImplCtx() (*BlscTypesImpl, errors.AIErrors) {
	return NewBlscTypeImpl("", "", "")
}

func NewBlscTypeImpl(endpoint string, accessKey string, accessSecretKey string) (*BlscTypesImpl, errors.AIErrors) {
	httpProfile, parseErr := ParseEndpoint(endpoint)
	if parseErr != errors.SUCCESS {
		fmt.Println(parseErr.Message)
		return nil, parseErr
	}
	return &BlscTypesImpl{
		HttpProfile: httpProfile,
		ClientProfile: common.ClientProfile{
			AccessKey:       accessKey,
			SecretAccessKey: accessSecretKey,
		},
	}, errors.SUCCESS
}

// DescribePublicImages 查询ACK产品的公共镜像信息。
// 返回DescribePublicImagesResponse，其中包含查询到的公共镜像信息，以及errors.AIErrors，用于处理请求过程中的错误。
// 本函数主要实现了与ACK产品服务端的通信，发送查询公共镜像的请求并处理响应。
func (prd *BlscTypesImpl) DescribePublicImages() (ack_product.DescribePublicImagesResponse, errors.AIErrors) {
	// 初始化ACK产品的客户端，用于发送HTTP请求。
	client := ack_product.ACKProductClient{
		HttpProfile:   prd.HttpProfile,
		ClientProfile: prd.ClientProfile,
	}

	// 创建DescribePublicImagesRequest对象，用于设置查询公共镜像的请求参数。
	request := ack_product.DescribePublicImagesRequest{}

	// 调用DescribeACKPublicImages方法查询公共镜像信息，返回查询结果和错误信息。
	describePublicImagesResponse, descErr := client.DescribeACKPublicImages(request)
	// 判断查询过程中是否发生错误，如果发生错误则返回空的DescribePublicImagesResponse和错误信息。
	if descErr != errors.SUCCESS {
		return ack_product.DescribePublicImagesResponse{}, descErr
	}
	// 如果查询成功，则返回查询结果和SUCCESS错误代码。
	return describePublicImagesResponse, errors.SUCCESS
}

// DescribeServiceTypes 查询ACK服务类型
// 本函数用于向ACK发送请求，以获取支持的服务类型列表。
// 返回值包括ACK服务类型的响应数据和可能的错误信息。
// 返回的错误信息使用自定义的错误类型AIErrors，以便更详细地描述错误情况。
func (prd *BlscTypesImpl) DescribeServiceTypes() (ack_product.DescribeServiceTypesResponse, errors.AIErrors) {
	// 初始化ACK产品的客户端
	// 使用当前实例的HTTP配置和客户端配置来创建ACK产品的客户端实例。
	client := ack_product.ACKProductClient{
		HttpProfile:   prd.HttpProfile,
		ClientProfile: prd.ClientProfile,
	}

	// 创建查询服务类型的请求对象
	request := ack_product.DescribeServiceTypesRequest{}

	// 调用客户端的DescribeACKServiceTypes方法查询服务类型
	// 并接收返回的响应和可能的错误
	describeServiceTypesResponse, descErr := client.DescribeACKServiceTypes(request)
	// 判断查询是否成功
	if descErr != errors.SUCCESS {
		// 如果查询失败，返回空的响应和错误信息
		return ack_product.DescribeServiceTypesResponse{}, descErr
	}

	// 如果查询成功，返回查询结果和成功标志
	return describeServiceTypesResponse, errors.SUCCESS
}

// DescribeSoldOut 用于查询指定区域和服务模型下的ACK产品售罄情况。
// 参数:
//
//	zoneCode - 可用区代码。
//	serviceModel - 服务模型。
//
// 返回值:
//
//	DescribeACKAvailableResourceResponse - 查询结果的响应对象。
//	errors.AIErrors - 错误码，成功时为errors.SUCCESS。
func (prd *BlscTypesImpl) DescribeSoldOut(zoneCode string, serviceModel string) (ack_product.DescribeACKAvailableResourceResponse, errors.AIErrors) {
	// 初始化ACK产品客户端
	client := ack_product.ACKProductClient{
		HttpProfile:   prd.HttpProfile,
		ClientProfile: prd.ClientProfile,
	}

	// 构造查询请求，指定服务模型和可用区
	request := ack_product.DescribeACKAvailableResourceRequest{
		ServiceModels: []common.ACKAvailableResourceInput{
			{
				ServiceModel: serviceModel,
				ZoneCode:     zoneCode,
			},
		},
	}

	// 发起查询请求，并处理可能的错误
	describeACKAvailableResourceResponse, descErr := client.DescribeACKAvailableResources(request)
	if descErr != errors.SUCCESS {
		return ack_product.DescribeACKAvailableResourceResponse{}, descErr
	}
	return describeACKAvailableResourceResponse, errors.SUCCESS
}
