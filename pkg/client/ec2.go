//go:generate mockgen -source=$GOFILE -destination=./ec2_mock.go -package=$GOPACKAGE -write_package_comment=false
package client

import (
	"context"
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Client interface {
	DescribeRegions(ctx context.Context) ([]string, error)
}

type EC2 struct {
	client *ec2.Client
}

var _ EC2Client = (*EC2)(nil)

func NewEC2(client *ec2.Client) *EC2 {
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
