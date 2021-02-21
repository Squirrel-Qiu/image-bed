package internal

import (
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Squirrel-Qiu/image-bed/store"
)

func (impl Implement) Get(ctx *gin.Context) {
	resourceId := ctx.Param("resourceId")

	bucket, ok, err := impl.DB.FileIsExistById(resourceId)
	if err != nil {
		logrus.Errorf("1 failed: %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "对象不存在"})
		return
	}

	cosClient := store.CloudClient{Credential: impl.Cred, Region: ""}
	cosConn := cosClient.Sign()

	getObjectOutput, err := cosConn.GetObject(&s3.GetObjectInput{
		Bucket:                     &bucket,
		ExpectedBucketOwner:        nil,
		IfMatch:                    nil,
		IfModifiedSince:            nil,
		IfNoneMatch:                nil,
		IfUnmodifiedSince:          nil,
		Key:                        &resourceId,
		PartNumber:                 nil,
		Range:                      nil,
		RequestPayer:               nil,
		ResponseCacheControl:       nil,
		ResponseContentDisposition: nil,
		ResponseContentEncoding:    nil,
		ResponseContentLanguage:    nil,
		ResponseContentType:        nil,
		ResponseExpires:            nil,
		SSECustomerAlgorithm:       nil,
		SSECustomerKey:             nil,
		SSECustomerKeyMD5:          nil,
		VersionId:                  nil,
	})
	if err != nil {
		logrus.Errorf("cosConn get object failed: %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "对象获取失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"": getObjectOutput.Body})
}
