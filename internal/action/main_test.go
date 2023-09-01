package action

import (
	"fmt"
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	fmt.Println()
	fmt.Println("==========================================")
	fmt.Println("========== Start Test: action ============")
	fmt.Println("==========================================")
	goleak.VerifyTestMain(m)
}
