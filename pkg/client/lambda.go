//go:generate mockgen -source=$GOFILE -destination=lambda_mock.go -package=$GOPACKAGE -write_package_comment=false
package client

import (
	"context"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type LambdaClient interface {
	ListFunctions(ctx context.Context) ([]types.FunctionConfiguration, error)
	ListFunctionsWithRegion(ctx context.Context, region string) ([]types.FunctionConfiguration, error)
	ListRuntimeValues() []string
}

type Lambda struct {
	client *lambda.Client
}

var _ LambdaClient = (*Lambda)(nil)

func NewLambda(client *lambda.Client) *Lambda {
	return &Lambda{
		client: client,
	}
}

func (c *Lambda) ListFunctions(ctx context.Context) ([]types.FunctionConfiguration, error) {
	return c.ListFunctionsWithRegion(ctx, "")
}

func (c *Lambda) ListFunctionsWithRegion(ctx context.Context, region string) ([]types.FunctionConfiguration, error) {
	var nextMarker *string
	outputs := []types.FunctionConfiguration{}

	var optFns func(*lambda.Options)
	if region != "" {
		optFns = func(o *lambda.Options) {
			o.Region = region
		}
	}

	for {
		input := &lambda.ListFunctionsInput{
			Marker: nextMarker,
		}

		var (
			output *lambda.ListFunctionsOutput
			err    error
		)

		if region == "" {
			output, err = c.client.ListFunctions(ctx, input)
		} else {
			output, err = c.client.ListFunctions(ctx, input, optFns)
		}
		if err != nil {
			return outputs, err
		}

		outputs = append(outputs, output.Functions...)

		nextMarker = output.NextMarker

		if nextMarker == nil {
			break
		}
	}

	return outputs, nil
}

func (c *Lambda) ListRuntimeValues() []string {
	var r types.Runtime
	runtimeStrList := []string{}
	runtimeList := r.Values()

	sort.Slice(runtimeList, func(i, j int) bool {
		first := string(runtimeList[i])
		second := string(runtimeList[j])

		firstRuntime, firstVersion, firstRest := c.splitVersion(first)
		secondRuntime, secondVersion, secondRest := c.splitVersion(second)

		if firstRuntime != secondRuntime {
			return firstRuntime < secondRuntime
		}

		if hasFinished, shouldSorted := c.compareActualVersion(firstVersion, secondVersion); hasFinished {
			return shouldSorted
		}

		if firstRest == "" {
			return true
		}
		if secondRest == "" {
			return false
		}

		return firstRest < secondRest
	})

	for _, runtime := range runtimeList {
		runtimeStrList = append(runtimeStrList, string(runtime))
	}

	return runtimeStrList
}

func (c *Lambda) splitVersion(runtimeStr string) (string, string, string) {
	r := regexp.MustCompile(`^(\D+)([\d\.]+)?(.*)?$`)
	matches := r.FindStringSubmatch(runtimeStr)

	runtime := ""
	version := ""
	rest := ""

	if len(matches) > 1 {
		runtime = matches[1]
	}
	if len(matches) > 2 {
		version = matches[2]
	}
	if len(matches) > 3 {
		rest = matches[3]
	}

	return runtime, version, rest
}

func (c *Lambda) compareActualVersion(first string, second string) (hasFinished bool, shouldSorted bool) {
	if first == "" {
		return true, true
	}
	if second == "" {
		return true, false
	}

	if first[:len(first)-1] == "." {
		first = first[len(first)-1:]
	}
	if second[:len(second)-1] == "." {
		second = second[len(second)-1:]
	}

	firstIntegers := first
	firstDecimals := ""
	secondIntegers := second
	secondDecimals := ""

	// test
	if i := strings.Index(first, "."); i >= 0 {
		firstIntegers = first[:i]
		firstDecimals = first[i+1:]
	}
	if i := strings.Index(second, "."); i >= 0 {
		secondIntegers = second[:i]
		secondDecimals = second[i+1:]
	}

	if firstIntegers != secondIntegers {
		fInt, _ := strconv.Atoi(firstIntegers)
		sInt, _ := strconv.Atoi(secondIntegers)
		return true, fInt < sInt
	}

	if firstDecimals == "" && secondDecimals != "" {
		return true, true
	}
	if firstDecimals != "" && secondDecimals == "" {
		return true, false
	}
	if firstDecimals != secondDecimals {
		fDec, _ := strconv.Atoi(firstDecimals)
		sDec, _ := strconv.Atoi(secondDecimals)
		return true, fDec < sDec
	}

	return false, false
}
