//go:generate mockgen -source=./lambda_factory.go -destination=./lambda_factory_mock.go -package=client
package client

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type LambdaCreator interface {
	Create(config aws.Config) *Lambda
}

type LambdaFactory struct{}

var _ LambdaCreator = (*LambdaFactory)(nil)

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
