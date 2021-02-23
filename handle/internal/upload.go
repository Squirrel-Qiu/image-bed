package internal

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Squirrel-Qiu/image-bed/store"
)

const URL = "https://..."

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

	bucket, resourceId, err := store.Store(hash, impl.DB)
	if err != nil {
		logrus.Errorf("db store failed: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// todo
	cosClient := store.CloudClient{Credential: impl.Cred, Region: ""}
	cosConn := cosClient.Sign()

	// check bucket if exist
	_, err = cosConn.HeadBucket(&s3.HeadBucketInput{Bucket: &bucket})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchBucket {
				logrus.Info("no such bucket, need to create")
				_, er := cosConn.CreateBucket(&s3.CreateBucketInput{Bucket: &bucket})
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
		Bucket:     &bucket,
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
