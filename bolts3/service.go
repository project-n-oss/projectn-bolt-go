package bolts3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/projectn-oss/projectn-bolt-go/boltv4"
	"os"
	"strings"
)

func New(p client.ConfigProvider, cfgs ...*aws.Config) *s3.S3 {
	boltUrl := os.Getenv("BOLT_URL")
	boltUrl = strings.Replace(boltUrl, "{region}", Region(p), -1)

	boltCfgs := aws.NewConfig().WithS3ForcePathStyle(true).WithEndpoint(boltUrl).WithRegion("us-east-1")
	boltCfgs.MergeIn(cfgs...)
	boltSvc := s3.New(p, boltCfgs)
	boltSvc.Handlers.Sign.Clear()
	boltSvc.Handlers.Sign.PushBackNamed(boltv4.SignBoltRequestHandler)

	return boltSvc
}

func Region(p client.ConfigProvider) string {

	region := os.Getenv("AWS_REGION")
	if len(region) > 0 {
		return region
	} else {
		ec2Md := ec2metadata.New(p)

		region, _ := ec2Md.Region()
		return region
	}
}
