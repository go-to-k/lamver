package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Ec2 struct {
	client *ec2.Client
}

func NewEC2Client(config aws.Config) *Ec2 {
	ec2Client := ec2.NewFromConfig(config)

	return &Ec2{
		client: ec2Client,
	}
}

func (ec2Client *Ec2) DescribeRegions(ctx context.Context) ([]string, error) {
	outputRegions := []string{}
	input := &ec2.DescribeRegionsInput{}

	output, err := ec2Client.client.DescribeRegions(ctx, input)
	if err != nil {
		return outputRegions, err
	}

	for _, region := range output.Regions {
		outputRegions = append(outputRegions, *region.RegionName)
	}

	return outputRegions, nil
}
