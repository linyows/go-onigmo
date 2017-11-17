package onigmo

import (
	"testing"
)

func TestOnigmoVersion(t *testing.T) {
	v := OnigmoVersion()
	if v != "6.1.3" {
		t.Errorf("OnigmoVersion wrong: %s", v)
	}
}
