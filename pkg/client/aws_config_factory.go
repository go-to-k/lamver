package client

import "context"

type AWSConfigCreator interface {
	Create(ctx context.Context, region string) (*AWSConfig, error)
}

type AWSConfigFactory struct {
	Profile string
}

var _ AWSConfigCreator = (*AWSConfigFactory)(nil)

func NewAWSConfigFactory(profile string) *AWSConfigFactory {
	return &AWSConfigFactory{
		Profile: profile,
	}
}

func (f *AWSConfigFactory) Create(ctx context.Context, region string) (*AWSConfig, error) {
	return NewAWSConfig(ctx, region, f.Profile)
}
