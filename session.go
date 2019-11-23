package main

import (
	"errors"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// JwtToken 은 CSI에서 사용하는 토큰 구조입니다.
type JwtToken struct {
	ID          string `json:"id"`
	LastProject string `json:"project"`
	AccessLevel `json:"accesslevel"`
	jwt.StandardClaims
}

// CreateTokenString 는 사용자의 기본 정보를 받아서 jwt token 키를 생성합니다.
func CreateTokenString(id string, accessLevel AccessLevel, lastProject string) (string, error) {
	// token에 정보를 넣는다. HS256 암호화 알고리즘을 사용합니다.
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &JwtToken{
		ID:          id,
		LastProject: lastProject,
		AccessLevel: accessLevel,
	})
	// 토큰에 CSI_JWT_SIGN_KEY를 가지고 와서 싸인합니다. 환경변수가 잡혀있지 않다면 빈 문자열로 싸인합니다.
	tokenstring, err := token.SignedString([]byte(os.Getenv("CSI_JWT_SIGN_KEY")))
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
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("CSI_JWT_SIGN_KEY")), nil
			})
			if err != nil {
				return jt, err
			}
			token, err = jwt.ParseWithClaims(cookie.Value, &jt, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("CSI_JWT_SIGN_KEY")), nil
			})
			if !token.Valid {
				return jt, errors.New("토큰이 유효하지 않습니다")
			}
			if err != nil {
				return jt, err
			}
			if jt.ID == "" {
				return jt, errors.New("ID가 빈 문자열입니다")
			}
			return jt, nil
		}
	}
	return jt, errors.New("token을 가지고 올 수 없습니다")
}

// RmSessionID 는 SessionID를 제거한다.
func RmSessionID(w http.ResponseWriter) {
	c := http.Cookie{
		Name:   "SSID",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, &c)
}
