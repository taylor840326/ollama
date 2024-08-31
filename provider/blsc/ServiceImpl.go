package blsc

import (
	"fmt"
	"github.com/taylor840326/aicloud-sdk-go/ackcs"
	"github.com/taylor840326/aicloud-sdk-go/common"
	"github.com/taylor840326/aicloud-sdk-go/common/errors"
	"github.com/taylor840326/aicloud-sdk-go/common/utils/StringUtils"
)

type BlscServiceImpl struct {
	HttpProfile   common.HttpProfile   `json:"httpProfile"`
	ClientProfile common.ClientProfile `json:"clientProfile"`
}

func NewBlscServiceImplCtx() (*BlscServiceImpl, errors.AIErrors) {
	return NewBlscServiceImpl("", "", "")
}

func NewBlscServiceImpl(endpoint string, accessKey string, accessSecretKey string) (*BlscServiceImpl, errors.AIErrors) {
	httpProfile, parseErr := ParseEndpoint(endpoint)
	if parseErr != errors.SUCCESS {
		fmt.Println(parseErr.Message)
		return nil, parseErr
	}
	return &BlscServiceImpl{
		HttpProfile: httpProfile,
		ClientProfile: common.ClientProfile{
			AccessKey:       accessKey,
			SecretAccessKey: accessSecretKey,
		},
	}, errors.SUCCESS
}

// CreateService 创建一个ACK服务。
// 本函数接收一个CreateServicesRequest类型的请求参数，返回一个CreateServicesResponse类型的响应数据和一个错误信息。
// 如果创建服务成功，错误信息为errors.SUCCESS；否则，返回相应的错误信息。
func (srv *BlscServiceImpl) CreateService(request ackcs.CreateServicesRequest) (ackcs.CreateServicesResponse, errors.AIErrors) {
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
// 返回errors.AIErrors表示操作是否成功。
func (srv *BlscServiceImpl) DeleteService(serviceUuid string) errors.AIErrors {
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
func (srv *BlscServiceImpl) DescribeServices(serviceUuid string) (ackcs.DescribeServicesResponse, errors.AIErrors) {
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
//	errors.AIErrors - 可能的错误代码和消息。
func (srv *BlscServiceImpl) DescribeServiceJupyter(serviceUuid string) (ackcs.DescribeServicesJupyterResponse, errors.AIErrors) {
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
//	errors.AIErrors - 可能的错误代码。
func (srv *BlscServiceImpl) DescribeServiceSSH(serviceUuid string) (ackcs.DescribeServicesSSHResponse, errors.AIErrors) {
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
