# Lambda Function Test Scripts

This directory contains test scripts for creating and deleting AWS Lambda functions.
These scripts are used for testing the `lamver` tool.

## Features

- Creates Lambda functions in multiple regions (us-east-1, ap-northeast-1)
- Supports multiple runtimes (Go, Node.js, Python)
- Automatically creates and deletes IAM roles and policies
- Allows specifying AWS profiles

## Usage

### Prerequisites

- AWS credentials configured (with appropriate IAM permissions)
- Go runtime environment

### Via Makefile

From the project root directory, you can use these commands:

```bash
# Default profile
make testgen

# Generate test Lambda functions with a specific profile
make testgen OPT="-p myprofile"

# Delete test Lambda functions
make testgen_clean

# Delete test Lambda functions with a specific profile
make testgen_clean OPT="-p myprofile"

# View help for test data generation
make testgen_help
```

### Direct Command Execution

You can also run the commands directly:

```bash
# Creating Lambda functions
cd testdata
go run cmd/create/main.go -p <profile>

# Deleting Lambda functions
cd testdata
go run cmd/delete/main.go -p <profile>
```

## Resource Naming

The test scripts create resources with the following naming patterns:

- IAM Role: `lamver-test-role`
- IAM Policy: `lamver-test-policy`
- Lambda Functions: `lamver-test-function-{region}-{runtime}`
  - Example: `lamver-test-function-us-east-1-nodejs`
