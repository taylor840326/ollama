package blsc

import (
	"github.com/taylor840326/aicloud-sdk-go/common"
	"github.com/taylor840326/aicloud-sdk-go/common/errors"
	"github.com/taylor840326/aicloud-sdk-go/common/utils/StringUtils"
	"regexp"
	"strconv"
)

const (
	RegexPattern = `(\w+)://(\w+.*):(\d+)`
)

// ParseEndpoint 解析给定的endpoint字符串，提取协议、主机、端口等信息。
// endpoint: 需要解析的endpoint字符串。
// 返回common.HttpProfile，包含解析后的协议、主机、端口等信息。
// 返回errors.AIErrors，可能包含解析过程中产生的错误。
func ParseEndpoint(endpoint string) (common.HttpProfile, errors.AIErrors) {
	// 检查endpoint是否为空，如果为空，则返回一个包含错误信息的HttpProfile和错误对象。
	if StringUtils.IsEmpty(endpoint) {
		return common.HttpProfile{}, errors.NewAIErrors(500, "endpoint is empty")
	}

	// 使用正则表达式编译提供的RegexPattern，用于后续的endpoint匹配。
	compile, regexErr := regexp.Compile(RegexPattern)
	if regexErr != nil {
		// 如果正则表达式编译失败，则返回一个包含编译错误信息的HttpProfile和错误对象。
		return common.HttpProfile{}, errors.NewAIErrorsWith(regexErr)
	}

	// 使用编译好的正则表达式对endpoint进行匹配，提取协议、主机、端口等信息。
	endpointParsed := compile.FindStringSubmatch(endpoint)

	// 将匹配结果中的端口号转换为整数类型。
	port, _ := strconv.Atoi(endpointParsed[3])

	// 返回包含解析后信息的HttpProfile和一个表示成功的错误对象。
	return common.HttpProfile{
		Protocol:       endpointParsed[1],
		Host:           endpointParsed[2],
		Port:           port,
		ContextPath:    "platform",
		ConnectTimeout: 1000,
		ReadTimeout:    1000,
	}, errors.SUCCESS
}
