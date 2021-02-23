package store

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
	Token     string
	Region    string
}

type CloudClient struct {
	Credential *Credential
	cosConn *s3.S3
}

func (me *CloudClient) Sign() *s3.S3 {
	if me.cosConn != nil {
		return me.cosConn
	}

	resolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.S3ServiceID {
			return endpoints.ResolvedEndpoint{
				URL:           fmt.Sprintf("https://cos.%s.myqcloud.com", region),
				SigningRegion: region,
			}, nil
		}
		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}

	creds := credentials.NewStaticCredentials(me.Credential.SecretId, me.Credential.SecretKey, me.Credential.Token)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      creds,
		Region:           aws.String(me.Credential.Region),
		EndpointResolver: endpoints.ResolverFunc(resolver),
	}))

	return s3.New(sess)
}
