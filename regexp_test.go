package onigmo

import (
	"fmt"
	"regexp"
	"testing"
)

func TestVersion(t *testing.T) {
	v := Version()
	if v != "6.1.3" {
		t.Errorf("Version wrong: %s", v)
	}
}

func TestSearchWithValidNamedGroup(t *testing.T) {
	s := "aaabbbbcc"
	regex, err := Compile("(?<foo>a*)(?<bar>b*)(?<foo>c*)")
	if err != nil {
		t.Error(err)
	}

	matched := regex.SearchString(s)
	if !matched {
		t.Error("Expected a match, but not a match")
	}

	foo, err := regex.matchResult.Get("foo")
	if err != nil {
		t.Error(err)
	}
	if foo != "aaa" {
		t.Errorf("Expected foo %v, but got %v", "aaa", foo)
	}

	bar, err := regex.matchResult.Get("bar")
	if err != nil {
		t.Error(err)
	}
	if bar != "bbbb" {
		t.Errorf("Expected bar %v, but got %v", "bbbb", bar)
	}

	defer regex.matchResult.Free()
	defer regex.Free()
}

func TestMatchWithValidNamedGroup(t *testing.T) {
	regex, err := Compile("^1st user (?<user>[a-z]*) ?2nd user (?<user>[a-z]+) value (?<val>[0-9]+)$")
	if err != nil {
		t.Error(err)
	}

	for _, data := range [][]string{
		[]string{"1st user foo 2nd user bar value 7", "foo", "7"},
		[]string{"1st user 2nd user bar value 789", "bar", "789"},
		[]string{"1st user somebody 2nd user else value 123", "somebody", "123"},
	} {
		matched := regex.MatchString(data[0])
		if !matched {
			t.Error("Expected a match")
		}

		user, err := regex.matchResult.Get("user")
		if err != nil {
			t.Error(err)
		}
		if user != data[1] {
			t.Errorf("Expected user %v, but got %v", data[1], user)
		}
		val, err := regex.matchResult.Get("val")
		if err != nil {
			t.Error(err)
		}
		if val != data[2] {
			t.Errorf("Expected val %v, but got %v", data[2], val)
		}
		defer regex.matchResult.Free()
	}

	defer regex.Free()
}

func TestMatchWithInValidNamedGroup(t *testing.T) {
	regex, err := Compile("^1st user (?<user>[a-z]*) ?2nd user (?<user>[a-z]+) (?<x>.*)(.*)value (?<val>[0-9]*)$")
	if err != nil {
		t.Error(err)
	}

	matched := regex.MatchString("1st user foo 2nd user bar value 789")
	if !matched {
		t.Error("Expected a match")
	}

	for _, data := range [][]string{
		[]string{"void", ""},
		[]string{"", ""},
	} {
		_, err := regex.matchResult.Get(data[0])
		if err == nil {
			t.Error("Expected error, because used non-existing capture group name.")
		}
	}

	val, err := regex.matchResult.Get("x")
	if err != nil {
		t.Error(err)
	}
	if val != "" {
		t.Errorf("Expected empty string, but got %v", val)
	}

	defer regex.matchResult.Free()
	defer regex.Free()
}

const STR = "1st user foo 2nd user bar value 7"

func BenchmarkOnigmo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re := MustCompile("^1st user (?<user>[a-z]*) ?2nd user (?<user>[a-z]+) value (?<val>[0-9]+)$")
		re.MatchString(STR)
		user, _ := re.matchResult.Get("user")
		if user != "foo" {
			fmt.Sprintf("Expected val %v, but got %v", "foo", user)
		}
		val, _ := re.matchResult.Get("val")
		if val != "7" {
			fmt.Sprintf("Expected val %v, but got %v", "7", val)
		}
		defer re.matchResult.Free()
		defer re.Free()
	}
}

func BenchmarkRE2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		re := regexp.MustCompile("^1st user ([a-z]*) ?2nd user ([a-z]+) value ([0-9]+)$")
		group := re.FindStringSubmatch(STR)
		if group[1] != "foo" {
			fmt.Sprintf("Expected user %v, but got %v", "foo", group[1])
		}
		if group[3] != "7" {
			fmt.Sprintf("Expected val %v, but got %v", "7", group[3])
		}
	}
}
