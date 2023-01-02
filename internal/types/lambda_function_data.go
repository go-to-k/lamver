package types

type LambdaFunctionData struct {
	Runtime      string
	Region       string
	FunctionName string
	LastModified string
}

func GetLambdaFunctionDataKeys() []string {
	return []string{"Runtime", "Region", "FunctionName", "LastModified"}
}
