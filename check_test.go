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

func TestTaskName(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{{
		in:   "comp", //소문자
		want: true,
	}, {
		in:   "comp1", //소문자, 숫자
		want: true,
	}, {
		in:   "comp_test", //소문자, 언더바
		want: true,
	}, {
		in:   "comp_1", //소문자, 언더바, 숫자
		want: true,
	}, {
		in:   "", //빈문자열
		want: false,
	}, {
		in:   "Comp", //대문자 포함
		want: false,
	}, {
		in:   "comp-1", //대쉬 포함
		want: false,
	}, {
		in:   "Comp-1", // 대문자, 대쉬 포함
		want: false,
	}}
	for _, c := range cases {
		got := regexpTask.MatchString(c.in)
		if got != c.want {
			t.Fatalf("TestShotname(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}

func TestStatusName(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{{
		in:   "wip", //소문자
		want: true,
	}, {
		in:   "pub1", //소문자, 숫자
		want: true,
	}, {
		in:   "wip_test", //소문자, 언더바
		want: false,
	}, {
		in:   "", //빈문자열
		want: false,
	}, {
		in:   "Wip", //대문자 포함
		want: true,
	}, {
		in:   "wip-1", //대쉬 포함
		want: false,
	}, {
		in:   "Wip-1", // 대문자, 대쉬 포함
		want: false,
	}}
	for _, c := range cases {
		got := regexpStatus.MatchString(c.in)
		if got != c.want {
			t.Fatalf("TestStatusName(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}

func TestEmail(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{{
		in:   "hello@lazypic.org",
		want: true,
	}, {
		in:   "hello@",
		want: false,
	}, {
		in:   "pm1@idea.co.kr",
		want: true,
	}, {
		in:   "d00000(test,test)",
		want: false,
	}}
	for _, c := range cases {
		got := regexpEmail.MatchString(c.in)
		if got != c.want {
			t.Fatalf("TestEmail(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}
