package blsc

import (
	"fmt"
	"github.com/taylor840326/aicloud-sdk-go/ackci"
	"github.com/taylor840326/aicloud-sdk-go/common"
	"github.com/taylor840326/aicloud-sdk-go/common/errors"
	"github.com/taylor840326/aicloud-sdk-go/common/utils/StringUtils"
)

type BlscImageImpl struct {
	HttpProfile   common.HttpProfile   `json:"httpProfile"`
	ClientProfile common.ClientProfile `json:"clientProfile"`
}

func NewBlscImageImplCtx() (*BlscImageImpl, errors.AIErrors) {
	return NewBlscImage("", "", "")
}
func NewBlscImage(endpoint string, accessKey string, accessSecretKey string) (*BlscImageImpl, errors.AIErrors) {
	httpProfile, parseErr := ParseEndpoint(endpoint)
	if parseErr != errors.SUCCESS {
		fmt.Println(parseErr.Message)
		return nil, parseErr
	}
	return &BlscImageImpl{
		HttpProfile: httpProfile,
		ClientProfile: common.ClientProfile{
			AccessKey:       accessKey,
			SecretAccessKey: accessSecretKey,
		},
	}, errors.SUCCESS
}

// CreateImages 根据请求创建镜像。
// 该方法封装了与ACK CI服务的交互，用于提交创建镜像的请求并处理响应。
//
// 参数:
//   - request: 包含创建镜像所需信息的请求对象。
//
// 返回值:
//   - CreateImagesResponse: 包含创建镜像操作的响应信息。
//   - AIErrors: 错误码，用于指示操作是否成功。
func (img *BlscImageImpl) CreateImages(request ackci.CreateImagesRequest) (ackci.CreateImagesResponse, errors.AIErrors) {
	// 初始化ACK CI客户端，使用当前实例的HTTP和客户端配置信息。
	client := ackci.ACKCIClient{
		HttpProfile:   img.HttpProfile,
		ClientProfile: img.ClientProfile,
	}

	// 调用客户端的CreateImages方法，提交创建镜像的请求。
	createImageResponse, createErr := client.CreateImages(request)
	// 检查创建镜像操作是否成功。
	if createErr != errors.SUCCESS {
		return ackci.CreateImagesResponse{}, createErr
	}
	// 如果操作成功，返回相应的响应信息和SUCCESS错误码。
	return createImageResponse, errors.SUCCESS
}

// DeleteImages 通过调用ACK容器服务接口，删除指定的镜像。
// 参数 imageUuid 为待删除镜像的UUID。
// 返回值为错误类型，如果删除成功，返回errors.SUCCESS。
func (img *BlscImageImpl) DeleteImages(imageUuid string) errors.AIErrors {
	// 初始化ACK客户端，使用img对象中的HttpProfile和ClientProfile配置信息。
	client := ackci.ACKCIClient{
		HttpProfile:   img.HttpProfile,
		ClientProfile: img.ClientProfile,
	}

	// 构造DeleteImages请求，指定区域代码和待删除的镜像UUID。
	request := ackci.DeleteImagesRequest{
		BaseRequest: common.BaseRequest{
			ZoneCode: "cn-zhongwei-ac",
		},
		ImageUuids: []string{imageUuid},
	}
	// 调用DeleteImages接口尝试删除镜像，返回响应和可能的错误。
	_, delErr := client.DeleteImages(request)
	// 如果删除操作出错，返回错误码；否则，返回errors.SUCCESS表示删除成功。
	if delErr != errors.SUCCESS {
		return delErr
	}
	return errors.SUCCESS
}

// DescribeImages 通过UUID查询镜像信息。
// 该方法用于向ACK CI服务发送请求，以获取指定镜像UUID的详细信息。
// 参数:
//
//	imageUuid - 需要查询的镜像的UUID。
//
// 返回值:
//
//	ackci.DescribeImagesResponse - 包含从ACK CI服务获取的镜像信息的响应。
//	errors.AIErrors - 如果查询过程中出现错误，则包含错误信息。
func (img *BlscImageImpl) DescribeImages(imageUuid string) (ackci.DescribeImagesResponse, errors.AIErrors) {
	// 初始化ACK CI客户端，使用img对象中的HttpProfile和ClientProfile配置。
	client := ackci.ACKCIClient{
		HttpProfile:   img.HttpProfile,
		ClientProfile: img.ClientProfile,
	}

	// 创建DescribeImagesRequest对象，用于向ACK CI服务发送查询请求。
	request := ackci.DescribeImagesRequest{
		PageDataRequest: common.PageDataRequest{
			PageNum:  1,
			PageSize: 10000,
		},
	}

	// 如果imageUuid不为空，则设置查询请求的ImageUuid字段。
	if !StringUtils.IsEmpty(imageUuid) {
		request.ImageUuid = imageUuid
	}

	// 调用ACK CI客户端的DescribeImages方法，发送查询请求，并获取响应及可能的错误。
	describeImagesResponse, descErr := client.DescribeImages(request)

	// 如果查询过程中出现错误，返回空的DescribeImagesResponse和错误信息。
	if descErr != errors.SUCCESS {
		return ackci.DescribeImagesResponse{}, descErr
	}

	// 如果查询成功，返回包含镜像信息的响应和SUCCESS错误代码。
	return describeImagesResponse, errors.SUCCESS
}
