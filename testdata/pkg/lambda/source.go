package lambda

import (
	"archive/zip"
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// Runtime information
type RuntimeInfo struct {
	Runtime  types.Runtime
	Source   string
	Handler  string
	FileExt  string
	FileName string
}

// Lambda function source code (Go)
const lambdaSourceGo = `
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Name string ` + "`json:\"name\"`" + `
}

type Response struct {
	Message string ` + "`json:\"message\"`" + `
}

func HandleRequest(ctx context.Context, event Event) (Response, error) {
	return Response{
		Message: fmt.Sprintf("Hello, %s!", event.Name),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
`

// Lambda function source code (Node.js)
const lambdaSourceNode = `
exports.handler = async (event) => {
    const name = event.name || 'World';
    return {
        message: "Hello, " + name + "!"
    };
};
`

// Lambda function source code (Python)
const lambdaSourcePython = `
def lambda_handler(event, context):
    name = event.get('name', 'World')
    return {
        'message': f'Hello, {name}!'
    }
`

// List of runtimes
var Runtimes = map[string]RuntimeInfo{
	"go": {
		Runtime:  types.RuntimeProvidedal2,
		Source:   lambdaSourceGo,
		Handler:  "main",
		FileExt:  ".go",
		FileName: "main.go",
	},
	"nodejs": {
		Runtime:  types.RuntimeNodejs22x,
		Source:   lambdaSourceNode,
		Handler:  "index.handler",
		FileExt:  ".js",
		FileName: "index.js",
	},
	"python": {
		Runtime:  types.RuntimePython313,
		Source:   lambdaSourcePython,
		Handler:  "lambda_function.lambda_handler",
		FileExt:  ".py",
		FileName: "lambda_function.py",
	},
}

// Compress Lambda function source code into a ZIP file
func createZip(runtimeInfo RuntimeInfo) ([]byte, error) {
	// Create ZIP file in memory
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Add source file to ZIP
	fileWriter, err := zipWriter.Create(runtimeInfo.FileName)
	if err != nil {
		return nil, err
	}

	_, err = io.WriteString(fileWriter, runtimeInfo.Source)
	if err != nil {
		return nil, err
	}

	// Close ZIP file
	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
