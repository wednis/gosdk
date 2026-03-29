package gosdk

// 身份验证相关

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrJWTInvalidMethod = errors.New("invalid method in jwt")
	ErrJWTInvalid       = errors.New("jwt invalid")
)

// 生成默认SigningMethodHS256（HMAC-SHA家族）的JWT
func NewJWT(data map[string]any, key string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, (jwt.MapClaims)(data)).SignedString([]byte(key))
}

// 验证JWT（属于HMAC-SHA家族且正确）
func VerifyJWT(s string, key string) (map[string]any, error) {
	token, err := jwt.Parse(s, func(token *jwt.Token) (any, error) {
		// 验证算法是否是HMAC-SHA家族内的（HS256就是）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrJWTInvalidMethod
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if jwtmap, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return jwtmap, nil
	}
	return nil, ErrJWTInvalid
}

// 生成Bcrypt加密
func NewBcrypt(source []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(source, 12)
}

// 验证Bcrypt加密
func VerifyBcrypt(source []byte, hash []byte) bool {
	if bcrypt.CompareHashAndPassword(hash, source) != nil {
		return true
	}
	return false
}

// 生成6位随机验证码
func NewCaptcha() []byte {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const length = 6
	code := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))
	for i := range length {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			panic(err)
		}
		code[i] = charset[num.Int64()]
	}
	return code
}
