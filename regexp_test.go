package onigmo

import (
	"testing"
)

func TestRgexp(t *testing.T) {
	v := Version()
	if v != "6.1.3" {
		t.Errorf("Version wrong: %s", v)
	}
}
