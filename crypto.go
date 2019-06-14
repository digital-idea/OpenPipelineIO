package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"math/rand"
	"strings"

	"gopkg.in/mgo.v2"
)

// Passcheck 함수는 사용 가능한 패스워드인지 체크하는 함수이다.
// MPAA의 기본조건은 7자리 이상 길이의 패스워드이다.
// 또한 MPAA 규칙중 10자리이상, 대문자, 소문자, 특수문자, 숫자, 연속된 문자금지 조항중 조건이 4개이상 일치해야한다.
func Passcheck(pw string) bool {
	if len(pw) <= 7 {
		return false
	}
	hasLong, hasUpper, hasLower, hasSpecial, hasNum, hasNotSame := false, false, false, false, false, true
	if len(pw) >= 10 { // 패스워드가 10자리 이상을 만족함.
		hasLong = true
	}
	var before rune
	for _, r := range pw {
		if before == r { // 전 문자열과 비교한다.
			hasNotSame = false
		}
		if strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZ", string(r)) {
			hasUpper = true
		}
		if strings.Contains("abcdefghijklmnopqrstuvwxyz", string(r)) {
			hasLower = true
		}
		if strings.Contains("0123456789", string(r)) {
			hasNum = true
		}
		if strings.Contains("~`!@#$%^&*()_-+={[}]|:;\"'<,>.?/", string(r)) {
			hasSpecial = true
		}
		before = r
	}
	condition := 0
	conditions := []bool{hasLong, hasUpper, hasLower, hasSpecial, hasNum, hasNotSame}
	for _, c := range conditions {
		if c {
			condition++
		}
	}
	if condition < 4 {
		return false
	}
	return true
}

// EncryptAES256 함수는 문자열을 AES256알고리즘으로 암호화한다.
func EncryptAES256(keystring, text string) (string, error) {
	key := []byte(KeygenAES256(keystring))
	/* block값의 예 :
		   블럭사이즈 내용은 암호화, 복호화값 입니다.
	       &{{[  808464432  808464432  808464432  808464432  808464432  808464432  808464432  808464432  875836469   67372037
	             875836469   67372037 3267543643 4076008043 3267543643 4076008043 1270726078 1337571771 2072874382 2139720075
	             275080550 3801585421  542463318 3534202685 1812521446  598825053 1480381907  666197080 3702533900 1042790401
	             510848343 3436530282 1849942584 1307753061  365850550  846935022 4287753252 3249951781 3754211698  319943448
	            3275644965 2395789888 2600646134 2843313688  740796553 3985767596  844449246  558151366 2009299787 4178551051
	            1645024509 3413414629  867559248 3726908412 3966852642 3442862308  511634832 3882962075 2239518822 1309198979 ]
	          [  511634832 3882962075 2239518822 1309198979  716807230 1450420625 1034901029  513984640 4151240949 4150400990
	            1592016643  984737963  987687460 2093626799 1809645492  588095141 2713479628     851755 2844141789 1683347368
	            4283349932 1175773067  387369499 1221681425 3403723630 2712704743 2844448758 3453222773 2276453005 3109708839
	            1359091088 1606920970  922975275 1800414601  138035473 1683564675 3630711175 1056280234 3898110391  247840410
	            1398428840 1548966306 1667916952 1818454418 3601727261 3869101869 3601727261 3869101869 1060715834  252251402
	            1060715834  252251402  808464432  808464432  808464432  808464432  808464432  808464432  808464432  808464432 ]}}
	*/
	block, err := aes.NewCipher(key) // 암호화를 위해서 새로운 암호화 블럭을 만듭니다.
	if err != nil {
		return "", err
	}
	// 문자열을 base64규칙(RFC4648)으로 바꿉니다.
	// base64는 ASCII영역의 문자들로만 이루어진 문자열입니다.
	base64str := base64.StdEncoding.EncodeToString([]byte(text))
	// 암호화 문자열을 저장할 변수 및 메모리영역을 생성합니다.
	// 만약 base64 값이 AAAA 라면 ciphertext 값은 [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 ] 이 됩니다. 16 + 4
	ciphertext := make([]byte, aes.BlockSize+len(base64str)) // aes.BlockSize값은 16입니다.
	// 일반적으로 AES형태의 암호화에서 iv값은 [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] 값이 됩니다.
	iv := ciphertext[:aes.BlockSize] // 맨앞의 16바이트는 암호화 스트림에 사용하기 위해서 iv변수로 빼놓습니다.

	// cfb 암호화 스트림을 만듭니다.
	// cfb 값 : &{0xc82025dad0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] 16 false}
	// 순서대로 dst, src, 암호화 블럭사이즈 입니다.
	cfb := cipher.NewCFBEncrypter(block, iv)

	// base64형태로 변환된 값을 이용해서 암호화에 사용되는 스트림(cfb)을 만듭니다.
	// cfb 값 : &{0xc82027a4e0 [215 40 160 42 200 73 68 55 46 222 218 194 68 1 165 212] [176 113 211 23 216 7 191 69 111 147 230 220 107 80 245 65] 4 false}
	// 여기까지 연산되면 ciphertext값은 아래의 형태를 띄게됩니다.
	// [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 77 230 10 223 101 62 232 231 ... 221 156 82 105]
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(base64str))
	return string(ciphertext[:len(ciphertext)]), nil
}

// DecryptAES256 함수는 AES256로 암호화된 문자열을 복호화한다.
func DecryptAES256(keystring, text string) (string, error) {
	key := []byte(KeygenAES256(keystring))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	bytestr := []byte(text)
	iv := bytestr[:aes.BlockSize]
	bytestr = bytestr[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(bytestr, bytestr)
	data, err := base64.StdEncoding.DecodeString(string(bytestr)) //RFC4648 방식의 인코딩
	if err != nil {
		return "", err
	}
	return string(data[:len(data)]), nil
}

// KeygenAES256 함수는 문자를 받아서 32자리 암호화를 위한키를 생성한다.
func KeygenAES256(str string) string {
	// 입력받은 문자열을 이용해서 렌덤 seed값을 먼저 만든다.
	// 문자 순서에 따라서 seed 값이 달라져야한다.
	// 따라서 매번 연산 부호를 바꾼다.
	var seed int64
	runes := []rune(str)
	var addOp bool
	for _, r := range runes {
		if addOp {
			seed += int64(r)
			addOp = false
			continue
		}
		seed *= int64(r)
		addOp = true
	}
	rand.Seed(seed)
	// seed값을 이용해서 64자리 AES256키 스팩의 랜덤문자를 생성한다.
	// aes키는 핸드폰에서도 추후 입력이 쉬워야한다.
	letters := []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	aeskey := make([]rune, 32)
	for i := range aeskey {
		aeskey[i] = letters[rand.Intn(len(letters))]
	}
	return string(aeskey)
}

// Validate 함수는 사용자의 ID,PW값이 유효한지 체크한다.
func Validate(id, pw string) bool {
	// id를 이용해서 해당 유저의 암호화된 pw를 가지고 온다.
	encryptPassword, err := EncryptAES256(pw, pw)
	if err != nil {
		return false
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		return false
	}
	defer session.Close()
	u, err := getUser(session, id)
	if u.Password != encryptPassword {
		return false
	}
	return true
}
