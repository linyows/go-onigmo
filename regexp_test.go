package onigmo

import (
	"testing"
)

func TestRgexp(t *testing.T) {
	v := OnigVersion()
	if v != "6.1.3" {
		t.Errorf("OnigVersion wrong: %s", v)
	}
}
