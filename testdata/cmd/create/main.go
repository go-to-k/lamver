package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/go-to-k/lamver/testdata/pkg/aws"
	"github.com/go-to-k/lamver/testdata/pkg/iam"
	"github.com/go-to-k/lamver/testdata/pkg/lambda"
)

func main() {
	// Parse command line arguments
	profile := flag.String("p", "", "AWS profile name")
	flag.Parse()

	// Store created function info
	var createdFunctions []string

	// Create IAM role in us-east-1 region first
	usEast1Cfg, err := aws.LoadConfig(context.TODO(), *profile, "us-east-1")
	if err != nil {
		log.Fatalf("Failed to load AWS configuration for us-east-1: %v", err)
	}

	roleARN, err := iam.CreateRole(context.TODO(), usEast1Cfg)
	if err != nil {
		log.Fatalf("Failed to create IAM role: %v", err)
	}

	fmt.Printf("Created/Reused IAM Role: %s\n", roleARN)

	// Wait for role propagation
	fmt.Println("Waiting for IAM role propagation (10 seconds)...")
	time.Sleep(10 * time.Second)

	// Create Lambda functions for each region and runtime combination
	for _, region := range aws.Regions {
		// Load AWS configuration for this region
		cfg, err := aws.LoadConfig(context.TODO(), *profile, region)
		if err != nil {
			log.Fatalf("Failed to load AWS configuration for %s: %v", region, err)
		}

		for runtimeName, runtimeInfo := range lambda.Runtimes {
			// Generate Lambda function name (including region and runtime)
			funcName := fmt.Sprintf("%s-%s-%s", lambda.FunctionName, region, runtimeName)

			// Create Lambda function
			isNew, err := lambda.CreateFunction(context.TODO(), cfg, funcName, roleARN, runtimeInfo)
			if err != nil {
				log.Fatalf("Failed to create Lambda function (%s, %s): %v", region, runtimeName, err)
			}

			// Different message based on whether the function already existed or was newly created
			if isNew {
				fmt.Printf("Lambda function '%s' successfully created in %s region (runtime: %s)\n",
					funcName, region, runtimeInfo.Runtime)
			} else {
				fmt.Printf("Lambda function '%s' already exists in %s region (runtime: %s) - reusing\n",
					funcName, region, runtimeInfo.Runtime)
			}

			// Store function info for display
			createdFunctions = append(createdFunctions, fmt.Sprintf("%s (Region: %s, Runtime: %s)",
				funcName, region, runtimeInfo.Runtime))
		}
	}

	// Display created functions summary
	fmt.Println("\nCreated Resources Summary:")
	fmt.Println("===========================")

	fmt.Println("\nIAM Resources:")
	fmt.Printf("- Role ARN: %s\n", roleARN)
	fmt.Printf("- Role Name: %s\n", iam.RoleName)
	fmt.Printf("- Policy Name: %s\n", iam.PolicyName)

	fmt.Println("\nLambda Functions:")
	for _, fn := range createdFunctions {
		fmt.Println("- " + fn)
	}
	fmt.Printf("\nTotal functions created: %d\n", len(createdFunctions))
}
