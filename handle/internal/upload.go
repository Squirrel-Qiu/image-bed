package internal

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/Squirrel-Qiu/image-bed/store"
)

func (impl Implement) Upload(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		logrus.Printf("upload file failed: %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "文件上传失败"})
		return
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("read file failed: %v", err)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(content))
	resourceId, ok := impl.DB.FileIsExistByHash(hash)
	if ok {
		logrus.Info("The upload file is exist")
		ctx.JSON(http.StatusOK, gin.H{"id": "url/" + resourceId})
		return
	}

	bucket, resourceId, err := store.Store(hash, impl.DB)
	if err != nil {
		logrus.Errorf("store resource failed: %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "文件存储失败"})
		return
	}

	cosClient := store.CloudClient{Credential: impl.Cred, Region: ""}
	cosConn := cosClient.Sign()

	// bucket is exist
	headBucketOutput, err := cosConn.HeadBucket(&s3.HeadBucketInput{Bucket: &bucket})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {

			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
				createBucketOutput, e := cosConn.CreateBucket(&s3.CreateBucketInput{
					ACL:                        nil,
					Bucket:                     &bucket,
					CreateBucketConfiguration:  nil,
					GrantFullControl:           nil,
					GrantRead:                  nil,
					GrantReadACP:               nil,
					GrantWrite:                 nil,
					GrantWriteACP:              nil,
					ObjectLockEnabledForBucket: nil,
				})
				if e != nil {
					logrus.Errorf("create bucket failed: %+v", e)
					ctx.JSON(http.StatusBadRequest, gin.H{"msg": "桶创建失败"})
					return
				}

			default:
				logrus.Errorf("create bucket failed: %+v", aerr.Error())
				ctx.JSON(http.StatusBadRequest, gin.H{"msg": "桶创建失败"})
				return
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and Message from an error.
			logrus.Errorf("create bucket failed: %+v", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "桶创建失败"})
			return
		}
	}

	putObjectOutput, err := cosConn.PutObject(&s3.PutObjectInput{
		ACL:                       nil,
		Body:                      file,
		Bucket:                    &bucket,
		BucketKeyEnabled:          nil,
		CacheControl:              nil,
		ContentDisposition:        nil,
		ContentEncoding:           nil,
		ContentLanguage:           nil,
		ContentLength:             nil,
		ContentMD5:                nil,
		ContentType:               nil,
		ExpectedBucketOwner:       nil,
		Expires:                   nil,
		GrantFullControl:          nil,
		GrantRead:                 nil,
		GrantReadACP:              nil,
		GrantWriteACP:             nil,
		Key:                       &resourceId,
		Metadata:                  nil,
		ObjectLockLegalHoldStatus: nil,
		ObjectLockMode:            nil,
		ObjectLockRetainUntilDate: nil,
		RequestPayer:              nil,
		SSECustomerAlgorithm:      nil,
		SSECustomerKey:            nil,
		SSECustomerKeyMD5:         nil,
		SSEKMSEncryptionContext:   nil,
		SSEKMSKeyId:               nil,
		ServerSideEncryption:      nil,
		StorageClass:              nil,
		Tagging:                   nil,
		WebsiteRedirectLocation:   nil,
	})
	if err != nil {
		logrus.Errorf("put object failed: %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "对象存储失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": "url/" + resourceId})
}

func postFile(resourceId string, file multipart.File) error {
	buff := new(bytes.Buffer)
	writer := multipart.NewWriter(buff)

	dst, err := writer.CreateFormFile("file", resourceId)
	if err != nil {
		return xerrors.Errorf("write to buffer: %w", err)
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		return xerrors.Errorf("copy file to cloud failed: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, "url", buff)
	if err != nil {
		return xerrors.Errorf("do request failed: %w", err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return xerrors.Errorf("do client failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
	}
	return nil
}
