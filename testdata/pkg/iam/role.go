package iam

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	awsConfig "github.com/go-to-k/lamver/testdata/pkg/aws"
)

const (
	RoleName         = "lamver-test-role"
	PolicyName       = "lamver-test-policy"
	ManagedPolicyARN = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
)

// IAM role policy document
const assumeRolePolicyDocument = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": {
				"Service": "lambda.amazonaws.com"
			},
			"Action": "sts:AssumeRole"
		}
	]
}`

// Lambda execution policy document
const lambdaExecutionPolicyDocument = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": [
				"logs:CreateLogGroup",
				"logs:CreateLogStream",
				"logs:PutLogEvents"
			],
			"Resource": "arn:aws:logs:*:*:*"
		}
	]
}`

// Create IAM role
func CreateRole(ctx context.Context, cfg aws.Config) (string, error) {
	client := iam.NewFromConfig(cfg)
	var roleArn string

	// Check if role already exists
	existingRole, err := client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String(RoleName),
	})

	// Role exists
	if err == nil && existingRole.Role != nil && existingRole.Role.Arn != nil {
		log.Printf("Role %s already exists, reusing it", RoleName)
		roleArn = *existingRole.Role.Arn
	} else {
		// Create role if it doesn't exist
		roleOutput, createErr := client.CreateRole(ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String(RoleName),
			AssumeRolePolicyDocument: aws.String(assumeRolePolicyDocument),
			Description:              aws.String("Role for lamver testing"),
		})
		if createErr != nil {
			return "", fmt.Errorf("failed to create role: %w", createErr)
		}
		roleArn = *roleOutput.Role.Arn
	}

	// Try to find existing policy
	accountID := getAccountID(ctx, cfg)
	customPolicyArn := fmt.Sprintf("arn:aws:iam::%s:policy/%s", accountID, PolicyName)

	// Try to get the policy
	_, policyErr := client.GetPolicy(ctx, &iam.GetPolicyInput{
		PolicyArn: aws.String(customPolicyArn),
	})

	// Create policy if it doesn't exist
	if policyErr != nil {
		policyOutput, createErr := client.CreatePolicy(ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(PolicyName),
			PolicyDocument: aws.String(lambdaExecutionPolicyDocument),
			Description:    aws.String("Policy for lamver testing"),
		})
		if createErr != nil {
			return "", fmt.Errorf("failed to create policy: %w", createErr)
		}
		customPolicyArn = *policyOutput.Policy.Arn
	} else {
		log.Printf("Policy %s already exists, reusing it", PolicyName)
	}

	// Check if policies are already attached
	attachedPolicies, err := client.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(RoleName),
	})

	var customPolicyAttached, awsManagedPolicyAttached bool
	if err == nil {
		for _, policy := range attachedPolicies.AttachedPolicies {
			if *policy.PolicyArn == customPolicyArn {
				customPolicyAttached = true
			} else if *policy.PolicyArn == ManagedPolicyARN {
				awsManagedPolicyAttached = true
			}
		}
	}

	// Attach custom policy if not already attached
	if !customPolicyAttached {
		_, err = client.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
			RoleName:  aws.String(RoleName),
			PolicyArn: aws.String(customPolicyArn),
		})
		if err != nil {
			return "", fmt.Errorf("failed to attach policy: %w", err)
		}
	}

	// Attach AWS managed policy if not already attached
	if !awsManagedPolicyAttached {
		_, err = client.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
			RoleName:  aws.String(RoleName),
			PolicyArn: aws.String(ManagedPolicyARN),
		})
		if err != nil {
			return "", fmt.Errorf("failed to attach AWS managed policy: %w", err)
		}
	}

	return roleArn, nil
}

// Get IAM role ARN
func GetRoleARN(ctx context.Context, cfg aws.Config) (string, error) {
	client := iam.NewFromConfig(cfg)

	output, err := client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String(RoleName),
	})
	if err != nil {
		return "", err
	}

	return *output.Role.Arn, nil
}

// Delete IAM roles and policies
func DeleteResources(ctx context.Context, profile string) error {
	// Delete IAM resources in us-east-1 region
	cfg, err := awsConfig.LoadConfig(ctx, profile, "us-east-1")
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	client := iam.NewFromConfig(cfg)
	accountID := getAccountID(ctx, cfg)

	// Detach custom policy
	_, err = client.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
		RoleName:  aws.String(RoleName),
		PolicyArn: aws.String(fmt.Sprintf("arn:aws:iam::%s:policy/%s", accountID, PolicyName)),
	})
	if err != nil {
		log.Printf("Warning: Failed to detach custom policy: %v", err)
	} else {
		fmt.Printf("Successfully detached policy '%s' from role '%s'\n", PolicyName, RoleName)
	}

	// Detach AWS managed policy
	_, err = client.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
		RoleName:  aws.String(RoleName),
		PolicyArn: aws.String(ManagedPolicyARN),
	})
	if err != nil {
		log.Printf("Warning: Failed to detach AWS managed policy: %v", err)
	} else {
		fmt.Printf("Successfully detached AWS managed policy from role '%s'\n", RoleName)
	}

	// Delete role
	_, err = client.DeleteRole(ctx, &iam.DeleteRoleInput{
		RoleName: aws.String(RoleName),
	})
	if err != nil {
		log.Printf("Warning: Failed to delete role: %v", err)
	} else {
		fmt.Printf("Successfully deleted IAM role '%s'\n", RoleName)
	}

	// Delete policy
	_, err = client.DeletePolicy(ctx, &iam.DeletePolicyInput{
		PolicyArn: aws.String(fmt.Sprintf("arn:aws:iam::%s:policy/%s", accountID, PolicyName)),
	})
	if err != nil {
		log.Printf("Warning: Failed to delete policy: %v", err)
	} else {
		fmt.Printf("Successfully deleted IAM policy '%s'\n", PolicyName)
	}

	return nil
}

// Get AWS account ID
func getAccountID(ctx context.Context, cfg aws.Config) string {
	client := iam.NewFromConfig(cfg)

	// Try GetUser
	result, err := client.GetUser(ctx, &iam.GetUserInput{})
	if err == nil && result.User != nil && result.User.Arn != nil {
		// Extract account ID from ARN
		arn := *result.User.Arn
		for i := 0; i < len(arn); i++ {
			if arn[i] == ':' {
				if i+1 < len(arn) && arn[i+1] == ':' {
					start := i + 2
					for j := start; j < len(arn); j++ {
						if arn[j] == ':' {
							return arn[start:j]
						}
					}
				}
			}
		}
	}

	// Try GetCallerIdentity
	stsClient := sts.NewFromConfig(cfg)
	identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err == nil && identity.Account != nil {
		return *identity.Account
	}

	log.Printf("Warning: Failed to get account ID: %v", err)
	return "*"
}
