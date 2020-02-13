package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestReverseStringSlice(t *testing.T) {
	cases := []struct {
		in   []string
		want []string
	}{{
		in:   nil,
		want: []string{},
	}, {
		in:   []string{},
		want: []string{},
	}, {
		in:   []string{"a"},
		want: []string{"a"},
	}, {
		in:   []string{"a", "b"},
		want: []string{"b", "a"},
	}}
	equal := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
	for _, c := range cases {
		got := ReverseStringSlice(c.in)
		if !equal(got, c.want) {
			t.Fatalf("RevereseStringSlice(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}

func TestShortTime(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{{
		in:   "",
		want: "",
	}, {
		in:   "1019",
		want: "1019",
	}, {
		in:   "2016-10-19",
		want: "1019",
	}, {
		in:   "2016-10-18T16:41:24+09:00",
		want: "1018",
	}, {
		in:   "2016-10-18T16:41:24-01:00",
		want: "1018",
	}, {
		in:   "2017-04-20T19:00:00+09:00",
		want: "0420",
	}}
	for _, c := range cases {
		got := ToShortTime(c.in)
		if ToShortTime(c.in) != c.want {
			t.Fatalf("ShortTime(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}

func TestFullTime(t *testing.T) {
	n := time.Now()
	// travisCI 에서는 UTC형식인 "2018-06-18T19:00:00Z" 라고 RFC3339 형식의 시간이 표기된다.
	timeOffset := n.Format(time.RFC3339)[19:len(n.Format(time.RFC3339))]
	cases := []struct {
		in   string
		want string
	}{{
		in:   "",
		want: "",
	}, {
		in:   "1019",
		want: fmt.Sprintf("%d-10-19T19:00:00%s", n.Year(), timeOffset),
	}, {
		in:   "2016-10-19",
		want: fmt.Sprintf("2016-10-19T19:00:00%s", timeOffset),
	}, {
		in:   "2016-10-18T16:41:24+09:00",
		want: "2016-10-18T16:41:24+09:00",
	}}
	for _, c := range cases {
		got := ToFullTime(c.in)
		if ToFullTime(c.in) != c.want {
			t.Fatalf("FullTime(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}

func TestCheckUpdate(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{{
		in:   "2016-10-28T18:42:25+09:00",
		want: false,
	}, {
		in:   "2016-10-28T18:42:25Z",
		want: false,
	}, {
		in:   time.Now().Format(time.RFC3339),
		want: true,
	}}
	for _, c := range cases {
		got := CheckUpdate(c.in)
		if CheckUpdate(c.in) != c.want {
			t.Fatalf("CheckUpdate(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}

func TestStr2Tags(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{{
		in:   "1,2,3",
		want: []string{"1", "2", "3"},
	}, {
		in:   " 1,2,3",
		want: []string{"1", "2", "3"},
	}, {
		in:   ",1,2,3",
		want: []string{"1", "2", "3"},
	}}
	for _, c := range cases {
		got := Str2Tags(c.in)
		if !reflect.DeepEqual(Str2Tags(c.in), c.want) {
			t.Fatalf("Str2Tags(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}

func TestList2str(t *testing.T) {
	cases := []struct {
		in   []string
		want string
	}{{
		in:   []string{"1", "2", "3"},
		want: "1,2,3",
	}, {
		in:   []string{"", "1", "2", "3"},
		want: "1,2,3",
	}}
	for _, c := range cases {
		got := List2str(c.in)
		if !reflect.DeepEqual(List2str(c.in), c.want) {
			t.Fatalf("Str2Tags(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}

func TestScanname2RollMedia(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{{
		in:   "",
		want: "",
	}, {
		in:   "22_A039C002_150916_R529",
		want: "A039C002",
	}, {
		in:   "2016-10-19",
		want: "",
	}, {
		in:   "A039C002_150916_R529",
		want: "",
	}}
	for _, c := range cases {
		got := Scanname2RollMedia(c.in)
		if got != c.want {
			t.Fatalf("FullTime(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}

func TestUsername2Elements(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{{
		in:   "",
		want: "",
	}, {
		in:   "김한웅(fire)",
		want: "fire",
	}, {
		in:   "김한웅(smoke),박지섭(fire)",
		want: "smoke,fire",
	}, {
		in:   "김한웅smoke),박상호(fire)",
		want: "fire",
	}}
	for _, c := range cases {
		got := Username2Elements(c.in)
		if got != c.want {
			t.Fatalf("FullTime(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}
