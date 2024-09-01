package blsc

import (
	"fmt"
	"github.com/taylor840326/aicloud-sdk-go/ack_product"
	"github.com/taylor840326/aicloud-sdk-go/ackci"
	"github.com/taylor840326/aicloud-sdk-go/ackcs"
	"github.com/taylor840326/aicloud-sdk-go/common"
	"github.com/taylor840326/aicloud-sdk-go/common/errors"
	"github.com/taylor840326/aicloud-sdk-go/common/utils/StringUtils"
)

type BlscImpl struct {
	HttpProfile   common.HttpProfile   `json:"httpProfile"`
	ClientProfile common.ClientProfile `json:"clientProfile"`
}

func NewBlscImplCtx() (*BlscImpl, error) {
	return NewBlscImpl("", "", "")
}

func NewBlscImpl(endpoint string, accessKey string, accessSecretKey string) (*BlscImpl, error) {
	httpProfile, parseErr := ParseEndpoint(endpoint)
	if parseErr != errors.SUCCESS {
		fmt.Println(parseErr.Message)
		return nil, parseErr
	}
	return &BlscImpl{
		HttpProfile: httpProfile,
		ClientProfile: common.ClientProfile{
			AccessKey:       accessKey,
			SecretAccessKey: accessSecretKey,
		},
	}, errors.SUCCESS
}

// DescribePublicImages 查询ACK产品的公共镜像信息。
// 返回DescribePublicImagesResponse，其中包含查询到的公共镜像信息，以及error，用于处理请求过程中的错误。
// 本函数主要实现了与ACK产品服务端的通信，发送查询公共镜像的请求并处理响应。
func (prd *BlscImpl) DescribePublicImages() (ack_product.DescribePublicImagesResponse, error) {
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
func (prd *BlscImpl) DescribeServiceTypes() (ack_product.DescribeServiceTypesResponse, error) {
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
//	error - 错误码，成功时为errors.SUCCESS。
func (prd *BlscImpl) DescribeSoldOut(zoneCode string, serviceModel string) (ack_product.DescribeACKAvailableResourceResponse, error) {
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

// CreateService 创建一个ACK服务。
// 本函数接收一个CreateServicesRequest类型的请求参数，返回一个CreateServicesResponse类型的响应数据和一个错误信息。
// 如果创建服务成功，错误信息为errors.SUCCESS；否则，返回相应的错误信息。
func (srv *BlscImpl) CreateService(request ackcs.CreateServicesRequest) (ackcs.CreateServicesResponse, error) {
	// 初始化ACKCSClient，使用srv中的HttpProfile和ClientProfile配置信息。
	client := ackcs.ACKCSClient{
		HttpProfile:   srv.HttpProfile,
		ClientProfile: srv.ClientProfile,
	}

	// 调用CreateServices方法创建服务，传入请求参数request。
	// 返回创建服务的响应信息和可能的错误。
	createServiceResponse, createErr := client.CreateServices(request)
	if createErr != errors.SUCCESS {
		// 如果创建服务失败，返回空的响应信息和错误信息。
		return ackcs.CreateServicesResponse{}, createErr
	}

	// 如果创建服务成功，返回创建服务的响应信息和errors.SUCCESS。
	return createServiceResponse, errors.SUCCESS
}

// DeleteService 通过服务UUID删除ACK服务。
// 返回error表示操作是否成功。
func (srv *BlscImpl) DeleteService(serviceUuid string) error {
	// 初始化ACKCS客户端
	client := ackcs.ACKCSClient{
		HttpProfile:   srv.HttpProfile,
		ClientProfile: srv.ClientProfile,
	}

	// 构造查询服务请求
	describeRequest := ackcs.DescribeServicesRequest{
		PageDataRequest: common.PageDataRequest{
			PageSize: 10000,
			PageNum:  1,
		},
		ServiceUuid: serviceUuid,
	}
	// 查询服务信息
	describeServicesResponse, descErr := client.DescribeServices(describeRequest)
	if descErr != errors.SUCCESS {
		return descErr
	}
	// 检查查询结果是否成功
	if describeServicesResponse.Code != common.AICloudResponseSuccessCode {
		return errors.NewAIErrors(500, describeServicesResponse.Message)
	}
	// 获取待删除的服务信息
	oneService := describeServicesResponse.Data.Rows[0]

	// 构造删除服务请求
	deleteRequest := ackcs.DeleteServicesRequest{
		BaseRequest: common.BaseRequest{
			ZoneCode: oneService.Zone.ZoneCode,
		},
		ServiceUuids: []string{serviceUuid},
	}
	// 删除服务
	deleteServicesResponse, delErr := client.DeleteServices(deleteRequest)
	if delErr != errors.SUCCESS {
		return delErr
	}
	// 检查删除操作是否成功
	if deleteServicesResponse.Code != common.AICloudResponseSuccessCode {
		return errors.NewAIErrors(500, deleteServicesResponse.Message)
	}

	return errors.SUCCESS
}

// DescribeServices 根据服务UUID查询ACK服务的详细信息。
// 服务UUID是查询的条件，如果为空，则返回所有服务的信息。
// 返回查询结果和可能的错误。
func (srv *BlscImpl) DescribeServices(serviceUuid string) (ackcs.DescribeServicesResponse, error) {
	// 初始化ACK客户端
	client := ackcs.ACKCSClient{
		HttpProfile:   srv.HttpProfile,
		ClientProfile: srv.ClientProfile,
	}

	// 创建查询服务的请求
	request := ackcs.DescribeServicesRequest{}
	// 设置每页返回1000条记录，从第1页开始
	request.PageSize = 1000
	request.PageNum = 1
	// 如果服务UUID不为空，则设置查询条件
	if !StringUtils.IsEmpty(serviceUuid) {
		request.ServiceUuid = serviceUuid
	}
	// 调用ACK服务查询接口
	describeResponse, descErr := client.DescribeServices(request)
	// 如果查询出错，则直接返回错误
	if descErr != errors.SUCCESS {
		return ackcs.DescribeServicesResponse{}, descErr
	}

	// 如果查询结果的返回码不是200，则构造错误并返回
	if describeResponse.Code != 200 {
		return ackcs.DescribeServicesResponse{}, errors.NewAIErrors(500, "describe services error "+describeResponse.Message)
	}

	// 返回查询结果和成功标志
	return describeResponse, errors.SUCCESS
}

// DescribeServiceJupyter 通过服务UUID查询Jupyter服务的详细信息。
// 该方法首先尝试获取服务的基本信息，然后根据这些信息获取特定Jupyter服务的详细配置。
// 参数:
//
//	serviceUuid - 服务的唯一标识符。
//
// 返回值:
//
//	ackcs.DescribeServicesJupyterResponse - 包含Jupyter服务详细信息的响应结构。
//	error - 可能的错误代码和消息。
func (srv *BlscImpl) DescribeServiceJupyter(serviceUuid string) (ackcs.DescribeServicesJupyterResponse, error) {
	// 初始化ACK CS客户端
	client := ackcs.ACKCSClient{
		HttpProfile:   srv.HttpProfile,
		ClientProfile: srv.ClientProfile,
	}

	// 查询服务基本信息
	descServicesResp, descErr := srv.DescribeServices(serviceUuid)
	if descErr != errors.SUCCESS {
		return ackcs.DescribeServicesJupyterResponse{}, descErr
	}
	// 检查查询服务的基本信息是否成功
	if descServicesResp.Code != common.AICloudResponseSuccessCode {
		return ackcs.DescribeServicesJupyterResponse{}, errors.NewAIErrors(500, "describe services jupyter error "+descServicesResp.Message)
	}

	// 获取服务列表
	services := descServicesResp.Data.Rows
	// 检查服务列表是否为空
	if len(services) == 0 {
		return ackcs.DescribeServicesJupyterResponse{}, errors.NewAIErrors(500, "describe services error ")
	}
	// 获取第一个服务实例
	oneService := services[0]

	// 构建查询Jupyter服务详情的请求
	request := ackcs.DescribeServicesJupyterRequest{
		BaseRequest: common.BaseRequest{
			ZoneCode: oneService.Zone.ZoneCode,
		},
		ServiceUuids: []string{serviceUuid},
	}

	// 发送查询Jupyter服务详情的请求
	servicesJupyterResponse, descErr := client.DescribeServicesJupyter(request)
	if descErr != errors.SUCCESS {
		return ackcs.DescribeServicesJupyterResponse{}, descErr
	}
	// 返回Jupyter服务的详细信息
	return servicesJupyterResponse, errors.SUCCESS
}

// DescribeServiceSSH 通过服务UUID获取服务的SSH信息。
// 该方法首先调用DescribeServices获取服务的详细信息，然后根据服务信息中的ZoneCode请求SSH信息。
// 参数:
//
//	serviceUuid - 服务的唯一标识符。
//
// 返回值:
//
//	ackcs.DescribeServicesSSHResponse - 包含SSH信息的响应。
//	error - 可能的错误代码。
func (srv *BlscImpl) DescribeServiceSSH(serviceUuid string) (ackcs.DescribeServicesSSHResponse, error) {
	// 初始化ACKCS客户端
	client := ackcs.ACKCSClient{
		HttpProfile:   srv.HttpProfile,
		ClientProfile: srv.ClientProfile,
	}

	// 调用DescribeServices获取服务详细信息
	descServicesResp, descErr := srv.DescribeServices(serviceUuid)
	if descErr != errors.SUCCESS {
		return ackcs.DescribeServicesSSHResponse{}, descErr
	}
	// 检查服务详细信息的响应码是否成功
	if descServicesResp.Code != common.AICloudResponseSuccessCode {
		return ackcs.DescribeServicesSSHResponse{}, errors.NewAIErrors(500, "describe services error "+descServicesResp.Message)
	}
	// 获取服务列表
	services := descServicesResp.Data.Rows
	// 检查服务列表是否为空
	if len(services) == 0 {
		return ackcs.DescribeServicesSSHResponse{}, errors.NewAIErrors(500, "describe services error ")
	}
	// 取出第一个服务
	oneService := services[0]

	// 构建DescribeServicesSSHRequest，请求SSH信息
	request := ackcs.DescribeServicesSSHRequest{
		BaseRequest: common.BaseRequest{
			ZoneCode: oneService.Zone.ZoneCode,
		},
		ServiceUuids: []string{serviceUuid},
	}

	// 调用DescribeServicesSSH获取SSH信息
	describeServicesSSHResponse, descErr := client.DescribeServicesSSH(request)
	if descErr != errors.SUCCESS {
		return ackcs.DescribeServicesSSHResponse{}, descErr
	}

	// 返回SSH信息的响应
	return describeServicesSSHResponse, errors.SUCCESS
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
func (img *BlscImpl) CreateImages(request ackci.CreateImagesRequest) (ackci.CreateImagesResponse, error) {
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
func (img *BlscImpl) DeleteImages(imageUuid string) error {
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
//	error - 如果查询过程中出现错误，则包含错误信息。
func (img *BlscImpl) DescribeImages(imageUuid string) (ackci.DescribeImagesResponse, error) {
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
