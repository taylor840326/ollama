package provider

import (
	"context"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/ollama/ollama/provider/blsc"
	"github.com/spf13/cobra"
	"github.com/taylor840326/aicloud-sdk-go/ackci"
	"github.com/taylor840326/aicloud-sdk-go/ackcs"
	"github.com/taylor840326/aicloud-sdk-go/common"
	"github.com/taylor840326/aicloud-sdk-go/common/enums"
	"github.com/taylor840326/aicloud-sdk-go/common/errors"
	"github.com/taylor840326/aicloud-sdk-go/common/utils/ArithUtils"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func DescribeServiceTypes(cmd *cobra.Command, args []string) error {

	typeService, newErr := blsc.NewBlscTypeImplCtx()
	if newErr != errors.SUCCESS {
		log.Fatalf(newErr.Message)
		return nil
	}

	describeServiceTypesResponse, descErr := typeService.DescribeServiceTypes()
	if descErr != errors.SUCCESS {
		log.Fatalf(descErr.Message)
		return nil
	}

	if describeServiceTypesResponse.Code != common.AICloudResponseSuccessCode {
		log.Fatalf(describeServiceTypesResponse.Message)
		return nil
	}
	serviceTypes := describeServiceTypesResponse.Data

	t := table.Table{}
	t.SetTitle("容器实例规格列表")
	header := table.Row{"可用区", "规格类型", "规格", "GPU卡数", "CPU核数", "内存大小", "是否售罄"}
	t.AppendHeader(header)
	for _, service := range serviceTypes {
		typeService, newErr := blsc.NewBlscTypeImplCtx()
		var soldOut bool = false
		if newErr == errors.SUCCESS {
			soldOutResponse, sdErr := typeService.DescribeSoldOut(service.Zone.ZoneCode, service.ServiceModel)
			if sdErr == errors.SUCCESS {
				soldOut = soldOutResponse.Data[0].SoldOut
			}
		}
		row := table.Row{service.Zone.ZoneCode, service.ServiceType, service.ServiceModel, service.ServiceGpus, service.ServiceCpus, ArithUtils.B2GB(int64(service.ServiceMemory)), soldOut}
		t.AppendRow(row)
	}
	fmt.Println(t.Render())
	return nil
}

func DescribeImages(cmd *cobra.Command, args []string) error {

	typeService, newErr := blsc.NewBlscTypeImplCtx()
	if newErr != errors.SUCCESS {
		log.Fatalf(newErr.Message)
		return nil
	}

	describePublicImagesResponse, descErr := typeService.DescribePublicImages()
	if descErr != errors.SUCCESS {
		log.Fatalf(descErr.Message)
		return nil
	}

	if describePublicImagesResponse.Code != common.AICloudResponseSuccessCode {
		log.Fatalf(describePublicImagesResponse.Message)
		return nil
	}
	publicImages := describePublicImagesResponse.Data

	t := table.Table{}
	t.SetTitle("容器镜像列表")
	header := table.Row{"镜像类型", "镜像ID", "操作系统", "镜像架构", "镜像架构版本", "镜像状态", "是否可删除"}
	t.AppendHeader(header)
	for _, publicImage := range publicImages {
		row := table.Row{"公共镜像", publicImage.ImageUuid, publicImage.ImageOs + " " + publicImage.ImageVersion, publicImage.PlatformName, publicImage.PlatformVersion, "Active", false}
		t.AppendRow(row)
	}

	pImageService, newErr := blsc.NewBlscImageImplCtx()
	if newErr != errors.SUCCESS {
		log.Fatalf(newErr.Message)
		return nil
	}

	privateImagesResponse, descErr := pImageService.DescribeImages("")
	if descErr != errors.SUCCESS {
		log.Fatalf(descErr.Message)
		return nil
	}
	if privateImagesResponse.Code != common.AICloudResponseSuccessCode {
		log.Fatalf(descErr.Message)
		return nil
	}
	privateImages := privateImagesResponse.Data.Rows
	for _, publicImage := range privateImages {
		row := table.Row{"自定义镜像", publicImage.ImageUuid, publicImage.ImageOs + " " + publicImage.ImageVersion, publicImage.PlatformName, publicImage.PlatformVersion, publicImage.ImageStatus, true}
		t.AppendRow(row)
	}
	fmt.Println(t.Render())
	return nil
}

func CreateServices(cmd *cobra.Command, args []string) error {

	zoneCode := ""
	aliasName := ""
	serviceModel := ""
	imageUuid := ""
	//zoneCode := context.Args().Get(0)
	//if StringUtils.IsEmpty(zoneCode) {
	//	return nil
	//}
	//aliasName := context.Args().Get(1)
	//if StringUtils.IsEmpty(aliasName) {
	//	return nil
	//}
	//serviceModel := context.Args().Get(2)
	//if StringUtils.IsEmpty(serviceModel) {
	//	return nil
	//}
	//imageUuid := context.Args().Get(3)
	//if StringUtils.IsEmpty(imageUuid) {
	//	return nil
	//}

	serviceService, newErr := blsc.NewBlscServiceImplCtx()
	if newErr != errors.SUCCESS {
		log.Fatalf(newErr.Message)
		return nil
	}

	request := ackcs.CreateServicesRequest{
		BaseRequest: common.BaseRequest{
			ZoneCode: zoneCode,
		},
		AliasName:    aliasName,
		ServiceModel: serviceModel,
		BillingType:  enums.PostPaid.String(),
		ImageUuid:    imageUuid,
		Count:        1,
	}
	createServicesResponse, createErr := serviceService.CreateService(request)
	if createErr != errors.SUCCESS {
		log.Fatalf(createErr.Message)
		return nil
	}
	log.Println(createServicesResponse.Message + ":" + createServicesResponse.Data[0].ResourceUuid)
	return nil
}

func DestroyServices(cmd *cobra.Command, args []string) error {

	serviceUuid := ""
	//serviceUuid := context.Args().Get(0)
	//if StringUtils.IsEmpty(serviceUuid) {
	//	log.Println("serviceUuid is empty")
	//	return nil
	//}

	serviceService, newErr := blsc.NewBlscServiceImplCtx()
	if newErr != errors.SUCCESS {
		log.Fatalf(newErr.Message)
		return nil
	}

	delErr := serviceService.DeleteService(serviceUuid)
	if delErr != errors.SUCCESS {
		log.Fatalf(delErr.Message)
		return nil
	}
	log.Println(serviceUuid + " delete successed")
	return nil
}

func CommitServicesChange(cmd *cobra.Command, args []string) error {

	serviceUuid := ""
	aliasName := ""
	//serviceUuid := context.Args().Get(0)
	//if StringUtils.IsEmpty(serviceUuid) {
	//	log.Println("serviceUuid is empty")
	//	return nil
	//}
	//aliasName := context.Args().Get(1)
	//if StringUtils.IsEmpty(aliasName) {
	//	log.Println("aliasName is empty")
	//	return nil
	//}

	serviceService, newErr := blsc.NewBlscServiceImplCtx()
	if newErr != errors.SUCCESS {
		log.Fatalf(newErr.Message)
		return nil
	}

	describeServicesResponse, descErr := serviceService.DescribeServices(serviceUuid)
	if descErr != errors.SUCCESS {
		log.Fatalf(descErr.Message)
		return nil
	}
	oneService := describeServicesResponse.Data.Rows[0]

	imageService, newErr := blsc.NewBlscImageImplCtx()
	if newErr != errors.SUCCESS {
		log.Fatalf(newErr.Message)
		return nil
	}

	createImageRequest := ackci.CreateImagesRequest{
		BaseRequest: common.BaseRequest{
			ZoneCode: oneService.Zone.ZoneCode,
		},
		ServiceUuid: serviceUuid,
		AliasName:   aliasName,
		ImageDesc:   "",
	}
	createImagesResponse, createErr := imageService.CreateImages(createImageRequest)
	if createErr != errors.SUCCESS {
		log.Fatalf(createErr.Message)
		return nil
	}
	imageUuid := createImagesResponse.Data[0].ResourceUuid
	log.Println(imageUuid + " create successed")
	return nil
}

func RemoveImages(cmd *cobra.Command, args []string) error {

	imageUuid := ""
	//imageUuid := context.Args().Get(0)
	//if StringUtils.IsEmpty(imageUuid) {
	//	log.Println("imageUuid is empty")
	//	return nil
	//}

	imageService, newErr := blsc.NewBlscImageImplCtx()
	if newErr != errors.SUCCESS {
		log.Fatalf(newErr.Message)
		return nil
	}

	createErr := imageService.DeleteImages(imageUuid)
	if createErr != errors.SUCCESS {
		log.Fatalf(createErr.Message)
		return nil
	}
	log.Println(imageUuid + " create successed")
	return nil
}

func DescribeService(cmd *cobra.Command, args []string) error {

	serviceService, newErr := blsc.NewBlscServiceImplCtx()
	if newErr != errors.SUCCESS {
		log.Fatalf(newErr.Message)
		return nil
	}
	services, descErr := serviceService.DescribeServices("")
	if descErr != errors.SUCCESS {
		log.Fatalf(descErr.Message)
		return nil
	}

	t := table.Table{}
	t.SetTitle("容器实例列表")
	header := table.Row{"可用区", "实例规格", "实例ID", "实例状态", "计费类型", "镜像ID"}
	t.AppendHeader(header)
	for _, service := range services.Data.Rows {

		row := table.Row{service.Zone.ZoneCode, service.InstanceType.ServiceModel, service.ServiceUuid, service.ServiceStatus, service.BillingType, service.Image.ImageUuid}
		t.AppendRow(row)
	}
	fmt.Println(t.Render())
	return nil
}

var commands = map[string][]string{
	"windows": []string{"cmd", "/c", "start"},
	"darwin":  []string{"open"},
	"linux":   []string{"xdg-open"},
}

func DescribeServiceJupyter(cmd *cobra.Command, args []string) error {

	// 获取命令行参数中的serviceUuid
	serviceUuid := ""
	//serviceUuid := context.Args().Get(0)
	//// 检查serviceUuid是否为空，为空则终止程序
	//if StringUtils.IsEmpty(serviceUuid) {
	//	log.Fatalf("请指定serviceUuid")
	//	return nil
	//}

	// 根据当前操作系统获取对应的运行命令
	run, ok := commands[runtime.GOOS]
	// 如果当前操作系统不支持，则终止程序
	if !ok {
		log.Fatalf("don't know how to open things on %s platform", runtime.GOOS)
		return nil
	}

	// 使用context初始化BlscServiceImplCtx
	serviceService, newErr := blsc.NewBlscServiceImplCtx()
	// 检查初始化是否成功，不成功则终止程序
	if newErr != errors.SUCCESS {
		log.Fatalf(newErr.Message)
		return nil
	}

	// 根据serviceUuid描述服务并获取Jupyter相关信息
	jupyters, descErr := serviceService.DescribeServiceJupyter(serviceUuid)
	// 检查描述服务是否成功，不成功则终止程序
	if descErr != errors.SUCCESS {
		log.Fatalf(descErr.Message)
		return nil
	}

	// 从返回的Jupyter信息中获取第一个Jupyter的URL
	url := jupyters.Data.Jupyters[0].Urls[0]
	// 创建并启动运行命令
	cmd2 := exec.Command(run[0], append(run[1:], url)...)
	return cmd2.Start()
}

func DescribeServiceSSH(cmd *cobra.Command, args []string) error {

	// 从命令行参数中获取serviceUuid
	serviceUuid := ""
	//serviceUuid := context.Args().Get(0)
	////// 检查serviceUuid是否为空
	////if StringUtils.IsEmpty(serviceUuid) {
	////	log.Fatalf("请指定serviceUuid")
	////	return nil
	//}

	// 初始化BlscServiceImplCtx，用于后续的服务查询
	serviceService, newErr := blsc.NewBlscServiceImplCtx()
	if newErr != errors.SUCCESS {
		log.Fatalf(newErr.Message)
		return nil
	}

	// 根据serviceUuid查询服务的SSH信息
	sshResp, sshErr := serviceService.DescribeServiceSSH(serviceUuid)
	if sshErr != errors.SUCCESS {
		fmt.Println(sshErr.Message)
		return nil
	}

	// 解析SSH连接地址
	url := strings.Split(sshResp.Data.Sshes[0].Url, "@")[1]

	// 使用SSH连接信息建立SSH客户端连接
	client, sshClientErr := ssh.Dial("tcp", url, &ssh.ClientConfig{
		User:            "pod",
		Auth:            []ssh.AuthMethod{ssh.Password(sshResp.Data.Sshes[0].Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if sshClientErr != nil {
		log.Fatalf("SSH dial error: %s", sshClientErr.Error())
		return nil
	}

	// 创建SSH会话
	session, sshSessErr := client.NewSession()
	if sshSessErr != nil {
		log.Fatalf("new session error: %s", sshSessErr.Error())
		return nil
	}
	// 确保会话在函数返回前关闭
	defer session.Close()

	// 将会话的输出和错误输出重定向到标准输出和错误输出
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// 设置终端模式
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	// 请求pty以支持交互式shell
	if sshSessionErr := session.RequestPty("linux", 32, 160, modes); sshSessionErr != nil {
		log.Fatalf("request pty error: %s", sshSessionErr.Error())
	}

	// 启动shell
	if sshSessionErr := session.Shell(); sshSessionErr != nil {
		log.Fatalf("start shell error: %s", sshSessionErr.Error())
	}

	// 等待shell会话结束
	if sshSessionErr := session.Wait(); sshSessionErr != nil {
		log.Fatalf("return error: %s", sshSessionErr.Error())
	}
	return nil
}

func CopyFile(cmd *cobra.Command, args []string) error {

	source := ""
	target := ""
	//source := context.Args().Get(0)
	//target := context.Args().Get(1)

	if !strings.Contains(source, ":") && !strings.Contains(target, ":") {
		log.Fatalf("请指定源文件路径和目标文件路径")
	}

	var serviceUuid string
	var locToRemote bool
	var sourcePath string
	var targetPath string
	if strings.Contains(source, ":") {
		serviceUuid = strings.Split(source, ":")[0]
		sourcePath = strings.Split(source, ":")[1]
		locToRemote = false
	} else {
		sourcePath = source
	}
	if strings.Contains(target, ":") {
		serviceUuid = strings.Split(target, ":")[0]
		targetPath = strings.Split(target, ":")[1]
		locToRemote = true
	} else {
		targetPath = target
	}
	log.Println(locToRemote)

	serviceService, newErr := blsc.NewBlscServiceImplCtx()
	if newErr != errors.SUCCESS {
		log.Fatalf(newErr.Message)
		return nil
	}

	sshResp, sshErr := serviceService.DescribeServiceSSH(serviceUuid)
	if sshErr != errors.SUCCESS {
		fmt.Println(sshErr.Message)
		return nil
	}

	url := strings.Split(sshResp.Data.Sshes[0].Url, "@")[1]

	clientConfig, _ := auth.PasswordKey("pod", sshResp.Data.Sshes[0].Password, ssh.InsecureIgnoreHostKey())

	client := scp.NewClient(url, &clientConfig)

	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return nil
	}

	f := &os.File{}
	if runtime.GOOS == "windows" {
		f, _ = os.Open("C:\\Windows\\System32\\drivers\\etc\\hosts")
	} else {
		f, _ = os.Open(sourcePath)
	}

	defer client.Close()

	defer f.Close()

	err = client.CopyFromFile(context.Background(), *f, targetPath, "0655")

	if err != nil {
		fmt.Println("Error while copying file ", err)
	}
	return nil
}

func versionHandler(cmd *cobra.Command, _ []string) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return
	}

	serverVersion, err := client.Version(cmd.Context())
	if err != nil {
		fmt.Println("Warning: could not connect to a running Ollama instance")
	}

	if serverVersion != "" {
		fmt.Printf("ollama version is %s\n", serverVersion)
	}

	if serverVersion != version.Version {
		fmt.Printf("Warning: client version is %s\n", version.Version)
	}
}

func NewProviderCmdCLI() *cobra.Command {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cobra.EnableCommandSorting = false

	if runtime.GOOS == "windows" {
		console.ConsoleFromFile(os.Stdin) //nolint:errcheck
	}

	rootCmd := &cobra.Command{
		Use:           "service",
		Short:         "Large language model runner",
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			if version, _ := cmd.Flags().GetBool("version"); version {
				versionHandler(cmd, args)
				return
			}

			cmd.Print(cmd.UsageString())
		},
	}

	rootCmd.Flags().BoolP("version", "v", false, "Show version information")

	createCmd := &cobra.Command{
		Use:     "create SERVICE",
		Short:   "Create a model from a Modelfile",
		Args:    cobra.ExactArgs(1),
		PreRunE: checkServerHeartbeat,
		RunE:    CreateHandler,
	}

	createCmd.Flags().StringP("file", "f", "Modelfile", "Name of the Modelfile")
	createCmd.Flags().StringP("quantize", "q", "", "Quantize model to this level (e.g. q4_0)")

	showCmd := &cobra.Command{
		Use:     "destroy SERVICE",
		Short:   "Show information for a model",
		Args:    cobra.ExactArgs(1),
		PreRunE: checkServerHeartbeat,
		RunE:    ShowHandler,
	}

	showCmd.Flags().Bool("license", false, "Show license of a model")
	showCmd.Flags().Bool("modelfile", false, "Show Modelfile of a model")
	showCmd.Flags().Bool("parameters", false, "Show parameters of a model")
	showCmd.Flags().Bool("template", false, "Show template of a model")
	showCmd.Flags().Bool("system", false, "Show system message of a model")

	listCmd := &cobra.Command{
		Use:     "ssh",
		Aliases: []string{"ls"},
		Short:   "ssh client",
		PreRunE: checkServerHeartbeat,
		RunE:    ListHandler,
	}

	copyCmd := &cobra.Command{
		Use:     "cp SOURCE DESTINATION",
		Short:   "Copy a model",
		Args:    cobra.ExactArgs(2),
		PreRunE: checkServerHeartbeat,
		RunE:    CopyHandler,
	}

	envVars := envconfig.AsMap()

	envs := []envconfig.EnvVar{envVars["OLLAMA_HOST"]}

	for _, cmd := range []*cobra.Command{
		createCmd,
		showCmd,
		runCmd,
		pullCmd,
		pushCmd,
		listCmd,
		psCmd,
		copyCmd,
		deleteCmd,
		serveCmd,
	} {
		switch cmd {
		case runCmd:
			appendEnvDocs(cmd, []envconfig.EnvVar{envVars["OLLAMA_HOST"], envVars["OLLAMA_NOHISTORY"]})
		case serveCmd:
			appendEnvDocs(cmd, []envconfig.EnvVar{
				envVars["OLLAMA_DEBUG"],
				envVars["OLLAMA_HOST"],
				envVars["OLLAMA_KEEP_ALIVE"],
				envVars["OLLAMA_MAX_LOADED_MODELS"],
				envVars["OLLAMA_MAX_QUEUE"],
				envVars["OLLAMA_MODELS"],
				envVars["OLLAMA_NUM_PARALLEL"],
				envVars["OLLAMA_NOPRUNE"],
				envVars["OLLAMA_ORIGINS"],
				envVars["OLLAMA_SCHED_SPREAD"],
				envVars["OLLAMA_TMPDIR"],
				envVars["OLLAMA_FLASH_ATTENTION"],
				envVars["OLLAMA_LLM_LIBRARY"],
			})
		default:
			appendEnvDocs(cmd, envs)
		}
	}

	rootCmd.AddCommand(
		serveCmd,
		createCmd,
		showCmd,
		runCmd,
		pullCmd,
		pushCmd,
		listCmd,
		psCmd,
		copyCmd,
		deleteCmd,
	)

	return rootCmd
}
