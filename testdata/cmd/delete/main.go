package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/go-to-k/lamver/testdata/pkg/aws"
	"github.com/go-to-k/lamver/testdata/pkg/iam"
	"github.com/go-to-k/lamver/testdata/pkg/lambda"
)

func main() {
	// Parse command line arguments
	profile := flag.String("p", "", "AWS profile name")
	flag.Parse()

	// Store deleted function info
	var deletedFunctions []string
	var failedFunctions []string

	// Delete Lambda functions for each region and runtime combination
	for _, region := range aws.Regions {
		// Load AWS configuration
		cfg, err := aws.LoadConfig(context.TODO(), *profile, region)
		if err != nil {
			log.Fatalf("Failed to load AWS configuration: %v", err)
		}

		// Delete Lambda functions for each runtime
		for runtimeName, runtimeInfo := range lambda.Runtimes {
			// Generate Lambda function name
			funcName := fmt.Sprintf("%s-%s-%s", lambda.FunctionName, region, runtimeName)

			// Delete Lambda function
			err = lambda.DeleteFunction(context.TODO(), cfg, funcName)
			if err != nil {
				log.Printf("Warning: Failed to delete Lambda function '%s': %v", funcName, err)
				failedFunctions = append(failedFunctions, fmt.Sprintf("%s (Region: %s, Runtime: %s) - Error: %v",
					funcName, region, runtimeInfo.Runtime, err))
				continue
			}

			fmt.Printf("Lambda function '%s' successfully deleted from %s region\n",
				funcName, region)

			// Store function info for display
			deletedFunctions = append(deletedFunctions, fmt.Sprintf("%s (Region: %s, Runtime: %s)",
				funcName, region, runtimeInfo.Runtime))
		}
	}

	// Delete IAM roles and policies (only in us-east-1)
	if err := iam.DeleteResources(context.TODO(), *profile); err != nil {
		log.Fatalf("Failed to delete IAM resources: %v", err)
	}

	// Display deleted resources summary
	fmt.Println("\nDeleted Resources Summary:")
	fmt.Println("===========================")

	fmt.Println("\nIAM Resources:")
	fmt.Printf("- Role Name: %s\n", iam.RoleName)
	fmt.Printf("- Policy Name: %s\n", iam.PolicyName)

	fmt.Println("\nLambda Functions:")
	for _, fn := range deletedFunctions {
		fmt.Println("- " + fn)
	}

	if len(failedFunctions) > 0 {
		fmt.Println("\nFailed to Delete:")
		for _, fn := range failedFunctions {
			fmt.Println("- " + fn)
		}
	}

	fmt.Printf("\nTotal functions deleted: %d\n", len(deletedFunctions))
	if len(failedFunctions) > 0 {
		fmt.Printf("Total failures: %d\n", len(failedFunctions))
	}
}
