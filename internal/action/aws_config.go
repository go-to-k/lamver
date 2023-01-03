package action

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

const DEFAULT_AWS_REGION = "us-east-1"

type AWSConfig struct {
	Config aws.Config
}

func NewAWSConfig(ctx context.Context, region string, profile string) (*AWSConfig, error) {
	var (
		err error
		cfg aws.Config
	)

	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return nil, err
	}

	if region != "" {
		cfg.Region = region
	}
	if cfg.Region == "" {
		cfg.Region = DEFAULT_AWS_REGION
	}

	return &AWSConfig{
		Config: cfg,
	}, nil
}
