package internal

import (
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Squirrel-Qiu/image-bed/store"
)

func (impl *Implement) Get(ctx *gin.Context) {
	resourceId := ctx.Param("resourceId")

	ok, bucket, err := impl.DB.FileIsExistById(resourceId)
	if !ok {
		if err != nil {
			logrus.Errorf("db check file is exist by id failed: %+v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		logrus.Info("the required file is not exist")
		ctx.Status(http.StatusBadRequest)
		return
	}

	cosClient := store.CloudClient{Credential: impl.Cred, Region: ""}
	cosConn := cosClient.Sign()

	getObjectOutput, err := cosConn.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &resourceId,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchKey {
				logrus.Errorf("the required file is not exist: %+v", err)
				ctx.Status(http.StatusBadRequest)
				return
			}
		} else {
			logrus.Errorf("cosConn get object failed: %+v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}

	if _, err := io.Copy(ctx.Writer, getObjectOutput.Body); err != nil {
		logrus.Errorf("copy object failed: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
}
