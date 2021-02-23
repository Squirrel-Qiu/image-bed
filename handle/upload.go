package handle

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Squirrel-Qiu/image-bed/store"
)

const (
	URL = "https://..."
	IdType = "resource_id"
)

type Resource struct {
	Id         string
	Bucket     string
	Hash       string
	CreateTime string
	Size       uint32
}

func (impl *Implement) Upload(ctx *gin.Context) {
	content, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logrus.Errorf("read file failed: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	h := sha256.Sum256(content)
	hash := hex.EncodeToString(h[:])
	ok, resourceId, err := impl.DB.FileIsExistByHash(hash)
	if !ok {
		if err != nil {
			logrus.Errorf("db check file is exist by hash failed: %+v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		logrus.Info("The upload file is exist")
		if _, err := ctx.Writer.WriteString(URL + resourceId); err != nil {
			logrus.Errorf("put object failed: %+v", err)
			return
		}
		return
	}

	resourceId, err = impl.Generator.GenerateId()
	if err != nil {
		logrus.Errorf("generate id failed: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	resource := new(Resource)
	resource.Id = resourceId
	now := time.Now()
	resource.Bucket = now.Format("2006-01")
	resource.CreateTime = now.Format("2006-01-02 15:04")
	resource.Hash = hash
	resource.Size = uint32(len(content))

	err = impl.DB.Store(resource)
	if err != nil {
		logrus.Errorf("db store resource failed: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	cosClient := store.CloudClient{Credential: impl.Cred}
	cosConn := cosClient.Sign()

	// check bucket if exist
	_, err = cosConn.HeadBucket(&s3.HeadBucketInput{Bucket: &resource.Bucket})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchBucket {
				logrus.Info("no such bucket, need to create")
				_, er := cosConn.CreateBucket(&s3.CreateBucketInput{Bucket: &resource.Bucket})
				if er != nil {
					logrus.Errorf("create bucket failed: %+v", er)
					ctx.Status(http.StatusInternalServerError)
					return
				}
			}
		} else {
			logrus.Errorf("create bucket failed: %+v", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}

	m := md5.Sum(content)
	cm := base64.StdEncoding.EncodeToString(m[:])
	reader := bytes.NewReader(content)
	_, err = cosConn.PutObject(&s3.PutObjectInput{
		Body:       reader,
		Bucket:     &resource.Bucket,
		ContentMD5: &cm,
		Key:        &resourceId,
	})
	if err != nil {
		logrus.Errorf("put object failed: %+v", err)
		ctx.Status(http.StatusInternalServerError)
	}

	if _, err := ctx.Writer.WriteString(URL + resourceId); err != nil {
		logrus.Errorf("put object failed: %+v", err)
		return
	}
}
