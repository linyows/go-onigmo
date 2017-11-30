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

func TestRegex(t *testing.T) {
	regex, err := Compile("^1st user (?<user>[a-z]*) ?2nd user (?<user>[a-z]+) value (?<value>[0-9]+)$")
	if err != nil {
		t.Errorf("Compile error: %#v", err)
	}

	result, err := regex.Match("1st user foo 2nd user bar value 5")
	if err != nil {
		t.Errorf("Match error: %#v", err)
	}

	user, err := result.Get("user")
	if err != nil {
		t.Errorf("Get error: %#v", err)
	}

	if user != "foo" {
		t.Errorf("User wrong: %s", user)
	}

	value, err := result.Get("value")
	if err != nil {
		t.Errorf("Get error: %#v", err)
	}

	if value != "5" {
		t.Errorf("Val wrong: %s", value)
	}

	defer regex.Free()
	defer result.Free()
}
