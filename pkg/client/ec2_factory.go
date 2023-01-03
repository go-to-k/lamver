package client

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func CreateEC2(config aws.Config) *EC2 {
	ec2Client := ec2.NewFromConfig(config, func(o *ec2.Options) {
		o.RetryMaxAttempts = retryMaxAttempts
		o.RetryMode = aws.RetryModeStandard
	})

	return NewEC2(ec2Client)
}
