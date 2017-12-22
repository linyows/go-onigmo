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

func TestValidCaptureGroups(t *testing.T) {
	regex, err := Compile("^1st user (?<user>[a-z]*) ?2nd user (?<user>[a-z]+) value (?<val>[0-9]+)$")
	if err != nil {
		t.Error(err)
	}

	for _, data := range [][]string{
		[]string{"1st user foo 2nd user bar value 7", "foo", "7"},
		// []string{"1st user 2nd user bar value 789", "bar", "789"},
		[]string{"1st user somebody 2nd user else value 123", "somebody", "123"},
	} {
		matched, err := regex.Match(data[0])
		if err != nil {
			t.Error(err)
		}
		if !matched {
			t.Error("expected a match")
		}

		user, err := regex.Get("user")
		if err != nil {
			t.Error(err)
		}
		if user != data[1] {
			t.Errorf("Expected user %v, but got %v", data[1], user)
		}
		val, err := regex.Get("val")
		if err != nil {
			t.Error(err)
		}
		if val != data[2] {
			t.Errorf("Expected val %v, but got %v", data[2], val)
		}
	}

	defer regex.Free()
}

func TestInvalidCaptureGroups(t *testing.T) {
	regex, err := Compile("^1st user (?<user>[a-z]*) ?2nd user (?<user>[a-z]+) (?<x>.*)(.*)value (?<val>[0-9]*)$")
	if err != nil {
		t.Error(err)
	}

	matched, err := regex.Match("1st user foo 2nd user bar value 789")
	if err != nil {
		t.Error(err)
	}
	if !matched {
		t.Error("expected a match")
	}

	for _, data := range [][]string{
		[]string{"void", ""},
		[]string{"", ""},
	} {
		_, err := regex.Get(data[0])
		if err == nil {
			t.Error("Expected error, because used non-existing capture group name.")
		}
	}
	val, err := regex.Get("x")
	if err != nil {
		t.Error(err)
	}
	if val != "" {
		t.Errorf("Expected empty string, but got %v", val)
	}

	defer regex.Free()
}
