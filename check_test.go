package main

import (
	"testing"
)

func TestShotname(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{{
		in:   "SS_0010",
		want: true,
	}, {
		in:   "",
		want: false,
	}, {
		in:   "SS_",
		want: false,
	}, {
		in:   "SS_0010_C001",
		want: false,
	}}
	for _, c := range cases {
		got := regexpShotname.MatchString(c.in)
		if got != c.want {
			t.Fatalf("FullTime(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}
