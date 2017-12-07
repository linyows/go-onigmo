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

	match, err := regex.Match("1st user foo 2nd user bar value 5")
	if err != nil {
		t.Errorf("Match error: %#v", err)
	}

	user, err := match.Get("user")
	if err != nil {
		t.Errorf("Get error: %#v", err)
	}

	if user != "foo" {
		t.Errorf("User wrong: %s", user)
	}

	value, err := match.Get("value")
	if err != nil {
		t.Errorf("Get error: %#v", err)
	}

	if value != "5" {
		t.Errorf("Val wrong: %s", value)
	}

	defer regex.Free()
	defer match.Free()
}

func TestInvalidCaptureGroups(t *testing.T) {
	regex, err := Compile("^1st user (?<user>[a-z]*) ?2nd user (?<user>[a-z]+) (?<x>.*)(.*)value (?<val>[0-9]*)$")
	if err != nil {
		t.Error(err)
	}

	match, err := regex.Match("1st user foo 2nd user bar value 789")
	if err != nil {
		t.Error(err)
	}

	if !match.IsMatch() {
		t.Error("expected a match")
	}
	for _, data := range [][]string{
		[]string{"void", ""},
		[]string{"", ""},
	} {
		_, err := match.Get(data[0])
		if err == nil {
			t.Error("Expected error, because used non-existing capture group name.")
		}
	}
	val, err := match.Get("x")
	if err != nil {
		t.Error(err)
	}
	if val != "" {
		t.Errorf("Expected empty string, but got %v", val)
	}

	defer regex.Free()
	defer match.Free()
}
