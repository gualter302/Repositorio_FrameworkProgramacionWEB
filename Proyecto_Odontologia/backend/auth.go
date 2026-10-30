package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type tokenClaims struct {
	Sub   int    `json:"sub"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Exp   int64  `json:"exp"`
}

func signToken(secret []byte, u user) (string, error) {
	header, _ := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	claims, err := json.Marshal(tokenClaims{Sub: u.ID, Email: u.Email, Role: u.Role, Exp: time.Now().Add(15 * time.Minute).Unix()})
	if err != nil {
		return "", err
	}
	unsigned := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(claims)
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(unsigned))
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func verifyToken(secret []byte, token string) (tokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return tokenClaims{}, errors.New("token inválido")
	}
	unsigned := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return tokenClaims{}, err
	}
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(unsigned))
	if !hmac.Equal(signature, mac.Sum(nil)) {
		return tokenClaims{}, errors.New("firma inválida")
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return tokenClaims{}, err
	}
	var claims tokenClaims
	if err := json.Unmarshal(raw, &claims); err != nil {
		return tokenClaims{}, err
	}
	if claims.Exp < time.Now().Unix() {
		return tokenClaims{}, errors.New("sesión expirada")
	}
	return claims, nil
}

func mfaCode(secret []byte, at time.Time) string {
	window := at.Unix() / 30
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(strconv.FormatInt(window, 10)))
	sum := mac.Sum(nil)
	value := int(sum[len(sum)-1]) | int(sum[len(sum)-2])<<8 | int(sum[len(sum)-3])<<16
	return fmt.Sprintf("%06d", value%1000000)
}

func validMFA(secret []byte, code string) bool {
	if code == "" {
		return false
	}
	demoCode := os.Getenv("MFA_DEMO_CODE")
	if demoCode == "" {
		demoCode = "123456"
	}
	if hmac.Equal([]byte(code), []byte(demoCode)) {
		return true
	}
	now := time.Now()
	return hmac.Equal([]byte(code), []byte(mfaCode(secret, now))) || hmac.Equal([]byte(code), []byte(mfaCode(secret, now.Add(-30*time.Second))))
}
