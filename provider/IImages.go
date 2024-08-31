package service

import (
	"github.com/taylor840326/aicloud-sdk-go/ackci"
	"github.com/taylor840326/aicloud-sdk-go/common/errors"
)

type IImages interface {
	CreateImages(request ackci.CreateImagesRequest) (ackci.CreateImagesResponse, errors.AIErrors)
	DeleteImages(imageUuid string) errors.AIErrors
	DescribeImages(imageUuid string) (ackci.DescribeImagesResponse, errors.AIErrors)
}
