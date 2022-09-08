package main

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// JwtToken 은 OpenPipelineIO에서 사용하는 토큰 구조입니다.
type JwtToken struct {
	jwt.StandardClaims        // 기본 Claims 구조 + 내가 필요한 Claim을 선언해서 사용한다.
	ID                 string `json:"id"`
	LastProject        string `json:"project"`
	AccessLevel        `json:"accesslevel"`
}

// CreateTokenString 는 사용자의 기본 정보를 받아서 jwt token 키를 생성합니다.
func CreateTokenString(id string, accessLevel AccessLevel, lastProject string) (string, error) {
	// token에 정보를 넣는다. HS256 암호화 알고리즘을 사용합니다.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":          id,
		"LastProject": lastProject,
		"AccessLevel": accessLevel,
	})
	// 토큰에 OPIO_JWT_SIGN_KEY를 가지고 서명키로 사용합니다. 환경변수가 잡혀있지 않다면 빈 문자열로 싸인합니다. 보안을 위해서 서명키 환경변수를 꼭 잡아주세요.
	tokenstring, err := token.SignedString([]byte(os.Getenv("OPIO_JWT_SIGN_KEY")))
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

// SetSessionID 는 SessionID를 가지고 온다.
func SetSessionID(w http.ResponseWriter, id string, accessLevel AccessLevel, project string) error {
	token, err := CreateTokenString(id, accessLevel, project)
	if err != nil {
		return err
	}
	c := http.Cookie{
		Name:    "SSID",
		Value:   token,
		Expires: time.Now().Add(time.Duration(*flagCookieAge) * time.Hour),
	}
	http.SetCookie(w, &c)
	return nil
}

// GetSessionID 는 SessionID를 가지고 온다.
func GetSessionID(r *http.Request) (JwtToken, error) {
	jt := JwtToken{}
	for _, cookie := range r.Cookies() {
		if cookie.Name == "SSID" {
			// 토큰이 파싱 가능한 형태인지 체크합니다.
			_, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(""), nil
			})
			if err != nil {
				return jt, err
			}
			// 파싱이 가능하하다는 것을 알았으니 파싱을 하고 파싱한 정보를 JWT 자료 구조에 넣습니다.
			token, err := jwt.ParseWithClaims(cookie.Value, &jt, func(token *jwt.Token) (interface{}, error) {
				return []byte(""), nil
			})
			if err != nil {
				return jt, err
			}
			if !token.Valid {
				return jt, errors.New("토큰이 유효하지 않습니다")
			}
			if jt.ID == "" {
				return jt, errors.New("ID가 빈 문자열입니다")
			}
			return jt, nil
		}
	}
	return jt, errors.New("Token 정보를 가지고 올 수 없습니다")
}

// RmSessionID 는 SessionID를 제거한다.
func RmSessionID(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "SSID",
		Value:  "",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "TOKEN",
		Value:  "",
		MaxAge: -1,
	})
}
