package action

import "context"

type AWSConfigFactory struct {
	Profile string
}

func NewAWSConfigFactory(profile string) *AWSConfigFactory {
	return &AWSConfigFactory{
		Profile: profile,
	}
}

func (f *AWSConfigFactory) Create(ctx context.Context, region string) (*AWSConfig, error) {
	return NewAWSConfig(ctx, region, f.Profile)
}
