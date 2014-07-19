package models

import (
	"fmt"
	"testing"
)


func TestTimeConsuming(t *testing.T) {
	if testing.Verbose() {
		// Error is equivalent to Log followed by Fail. Fail marks the function as having failed but continues execution.
		t.Error("fail message")
		t.Skip("skipping test in short mode.")
	}
}

func TestOther(t *testing.T) {
	//t.Fatal("fatal error, end this test now!")
}

func ExampleModels() {
	fmt.Println("hello")
	// Output: hello
}

func ExampleLoadModelsConfig() {
	fmt.Println("hello, and")
	fmt.Println("goodbye")
	// Output:
	// hello, and
	// goodbye
}
