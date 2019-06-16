package main

import (
	"testing"
)

func Test_Passcheck(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{{
		in:   "Aa!b89ah",
		want: true,
	}, {
		in:   "idea",
		want: false,
	}, {
		in:   "1234567890",
		want: false,
	}, {
		in:   "asdflk234",
		want: false,
	}, {
		in:   "asdf!lk233", // 10자 이상이라서 OK
		want: true,
	}, {
		in:   "af!lk233",
		want: false,
	}}
	for _, c := range cases {
		got := Passcheck(c.in)
		if Passcheck(c.in) != c.want {
			t.Fatalf("ShortTime(%v): 얻은 값 %v, 원하는 값 %v", c.in, got, c.want)
		}
	}
}
