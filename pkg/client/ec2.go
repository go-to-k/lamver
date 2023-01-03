package client

import (
	"context"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2 struct {
	client *ec2.Client
}

func NewEC2Client(config aws.Config) *EC2 {
	ec2Client := ec2.NewFromConfig(config, func(o *ec2.Options) {
		o.RetryMaxAttempts = retryMaxAttempts
		o.RetryMode = aws.RetryModeStandard
	})

	return &EC2{
		client: ec2Client,
	}
}

func (c *EC2) DescribeRegions(ctx context.Context) ([]string, error) {
	outputRegions := []string{}
	input := &ec2.DescribeRegionsInput{}

	output, err := c.client.DescribeRegions(ctx, input)
	if err != nil {
		return outputRegions, err
	}

	for _, region := range output.Regions {
		outputRegions = append(outputRegions, *region.RegionName)
	}

	sort.Strings(outputRegions)
	return outputRegions, nil
}
