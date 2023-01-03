package client

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type LambdaFactory struct{}

func NewLambdaFactory() *LambdaFactory {
	return &LambdaFactory{}
}

func (f *LambdaFactory) Create(config aws.Config) *Lambda {
	lambdaClient := lambda.NewFromConfig(config, func(o *lambda.Options) {
		o.RetryMaxAttempts = retryMaxAttempts
		o.RetryMode = aws.RetryModeStandard
	})

	return NewLambda(lambdaClient)
}
