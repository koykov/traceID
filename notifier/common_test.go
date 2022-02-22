package notifier

import "testing"

func assertStr(t testing.TB, n, a, b string) {
	if a != b {
		t.Errorf("%s mismatch: need '%s', got '%s'", n, a, b)
	}
}
