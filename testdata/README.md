# Lambda Function Test Scripts

This directory contains test scripts for creating and deleting AWS Lambda functions.
These scripts are used for testing the lamver tool.

## Features

- Creates Lambda functions in multiple regions (us-east-1, ap-northeast-1)
- Supports multiple runtimes (Go, Node.js, Python)
- Automatically creates and deletes IAM roles and policies
- Allows specifying AWS profiles

## Usage

### Prerequisites
- AWS credentials configured (with appropriate IAM permissions)
- Go runtime environment

### Command Options

The following options are supported:

- `-p, --profile`: AWS profile name
- Additional options defined in `cmd/create/main.go` and `cmd/delete/main.go`

### Via Makefile

From the project root directory, you can use these commands:

```bash
# Generate test Lambda functions
make testgen OPT="-p myprofile"

# Delete test Lambda functions
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

## Implementation Details

- Uses AWS SDK for Go v2 to create and delete Lambda functions
- Lambda functions are implemented with simple inline code
- Creates functions with various naming patterns to test filtering capabilities
- Uses different runtimes to test runtime filtering features
- Automatically handles existing resources (reuses instead of failing)
- Created functions can be used to test the lamver tool's region and runtime filtering

## Resource Naming

The test scripts create resources with the following naming patterns:

- IAM Role: `lamver-test-role`
- IAM Policy: `lamver-test-policy`
- Lambda Functions: `lamver-test-function-{region}-{runtime}`
  - Example: `lamver-test-function-us-east-1-nodejs`
