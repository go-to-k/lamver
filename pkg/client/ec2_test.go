package client

import (
	"context"
	"reflect"
	"testing"
)

func TestNewEC2(t *testing.T) {
	type args struct {
		client EC2SDKClient
	}
	tests := []struct {
		name string
		args args
		want *EC2
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEC2(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEC2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEC2_DescribeRegions(t *testing.T) {
	type fields struct {
		client EC2SDKClient
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EC2{
				client: tt.fields.client,
			}
			got, err := c.DescribeRegions(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("EC2.DescribeRegions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EC2.DescribeRegions() = %v, want %v", got, tt.want)
			}
		})
	}
}
