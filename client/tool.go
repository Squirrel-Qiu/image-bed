package client

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type Tool interface {
	Storage(resourceId, bucket string, reader *bytes.Reader) error
	Take(resourceId, bucket string) (reader io.Reader, err error)
}

func (c *Conn) Storage(resourceId, bucket string, reader *bytes.Reader) error {
	bucket = fmt.Sprintf("%s-%s", bucket, c.APPID)
	// check bucket if exist
	_, err := c.CosConn.HeadBucket(&s3.HeadBucketInput{Bucket: &bucket})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {

			if aerr.Code() == "NotFound" {
				logrus.Info("no such bucket, need to create")
				_, er := c.CosConn.CreateBucket(&s3.CreateBucketInput{Bucket: &bucket})
				if er != nil {
					return xerrors.Errorf("create bucket failed: %w", er)
				}
			}

		} else {
			return xerrors.Errorf("create bucket failed: %w", err)
		}
	}

	_, err = c.CosConn.PutObject(&s3.PutObjectInput{
		Body:   reader,
		Bucket: &bucket,
		Key:    &resourceId,
	})
	if err != nil {
		return xerrors.Errorf("put object failed: %w", err)
	}

	return nil
}

func (c *Conn) Take(resourceId, bucket string) (reader io.Reader, err error) {
	bucket = fmt.Sprintf("%s-%s", bucket, c.APPID)
	getObjectOutput, err := c.CosConn.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &resourceId,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchKey {
				return nil, xerrors.Errorf("the required file is not exist: %w", err)
			}
		} else {
			return nil, xerrors.Errorf("cosConn get object failed: %w", err)
		}
	}

	return getObjectOutput.Body, nil
}
