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
		in:   "SS-0010",
		want: true,
	}, {
		in:   "ss_0010",
		want: true,
	}, {
		in:   "ss0010",
		want: true,
	}, {
		in:   "",
		want: false,
	}, {
		in:   "SS_",
		want: false,
	}, {
		in:   "SS_0010_C001",
		want: true,
	}, {
		in:   "SS_00100",
		want: true,
	}, {
		in:   "SS-00100",
		want: true,
	}, {
		in:   "S001_C0010",
		want: true,
	}, {
		in:   "S001-C0010",
		want: true,
	}, {
		in:   "S001_C0010A",
		want: true,
	}, {
		in:   "001_0010",
		want: true,
	}, {
		in:   "S001-C0010A",
		want: true,
	}, {
		in:   "S001_C0010_A",
		want: true,
	}}
	for _, c := range cases {
		got := regexpShotname.MatchString(c.in)
		if got != c.want {
			t.Fatalf("TestShotname(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}

func TestVaildRnumTags(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{{
		in:   "1권",
		want: true,
	}, {
		in:   "24권",
		want: true,
	}, {
		in:   "SS_",
		want: false,
	}, {
		in:   "권",
		want: false,
	}}
	for _, c := range cases {
		got := validRnumTag(c.in)
		if got != c.want {
			t.Fatalf("TestValiedRnumTags(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}
