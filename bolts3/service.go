package bolts3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/projectn-oss/projectn-bolt-go/boltv4"
	"os"
)

func New(p client.ConfigProvider, cfgs ...*aws.Config) *s3.S3 {
	boltUrl := os.Getenv("BOLT_URL")

	boltCfgs := aws.NewConfig().WithS3ForcePathStyle(true).WithEndpoint(boltUrl).WithRegion("us-east-1")
	boltCfgs.MergeIn(cfgs...)
	boltSvc := s3.New(p, boltCfgs)
	boltSvc.Handlers.Sign.Clear()
	boltSvc.Handlers.Sign.PushBackNamed(boltv4.SignBoltRequestHandler)

	return boltSvc
}
