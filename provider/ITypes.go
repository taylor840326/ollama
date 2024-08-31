package service

import (
	"github.com/taylor840326/aicloud-sdk-go/ack_product"
	"github.com/taylor840326/aicloud-sdk-go/common/errors"
)

type ITypes interface {
	DescribeServiceTypes() (ack_product.DescribeServiceTypesResponse, errors.AIErrors)
	DescribePublicImages() (ack_product.DescribePublicImagesResponse, errors.AIErrors)
	DescribeSoldOut() (ack_product.DescribeACKAvailableResourceResponse, errors.AIErrors)
}
