package main

import (
	"fmt"
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

//테스트문자열을 암호화, 복호화 하기 위해서 테스트한 코드이다.
func Test_AES256(t *testing.T) {
	testString := "김한웅 woong 漢雄 か" // 한국어, 영어, 중국어, 일본어를 넣고 테스트
	e, err := EncryptAES256("!2A#adff", testString)
	if err != nil {
		fmt.Println(err)
	}
	d, err := DecryptAES256("!2A#adff", e)
	if err != nil {
		fmt.Println(err)
	}
	if d != testString {
		t.Fatalf("Test_AES256 : 입력값 : %v, 암호화값 %v, 복호화값 %v", testString, e, d)
	}
}
