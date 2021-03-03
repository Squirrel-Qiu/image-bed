package client

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Credential struct {
	SecretId  string
	SecretKey string
	APPID     string
	Region    string
}

type Conn struct {
	CosConn *s3.S3
	APPID   string
}

func Sign(c *Credential) *Conn {
	resolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.S3ServiceID {
			return endpoints.ResolvedEndpoint{
				URL:           fmt.Sprintf("https://cos.%s.myqcloud.com", region),
				SigningRegion: region,
			}, nil
		}
		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}

	creds := credentials.NewStaticCredentials(c.SecretId, c.SecretKey, "")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      creds,
		Region:           aws.String(c.Region),
		EndpointResolver: endpoints.ResolverFunc(resolver),
	}))

	return &Conn{CosConn: s3.New(sess), APPID: c.APPID}
}
