//go:generate mockgen -source=./ec2.go -destination=./ec2_mock.go -package=client
package client

import (
	"context"
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Client interface {
	DescribeRegions(ctx context.Context) ([]string, error)
}

type EC2SDKClient interface {
	DescribeRegions(ctx context.Context, params *ec2.DescribeRegionsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRegionsOutput, error)
}

type EC2 struct {
	client EC2SDKClient
}

var _ EC2Client = (*EC2)(nil)

func NewEC2(client EC2SDKClient) *EC2 {
	return &EC2{
		client: client,
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
