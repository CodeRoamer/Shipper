package routers

import (
	"testing"
)

func TestTimeConsuming(t *testing.T) {
	if testing.Verbose() {
		// Error is equivalent to Log followed by Fail. Fail marks the function as having failed but continues execution.
		t.Error("fail message")
		t.Skip("skipping test in short mode.")
	}
}
