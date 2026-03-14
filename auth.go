package gosdk

// 身份验证相关

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

var (
	ErrJWTInvaildMethod = errors.New("invaild method in jwt")
	ErrJWTInvaild       = errors.New("jwt invaild")
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
			return nil, ErrJWTInvaildMethod
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if jwtmap, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return jwtmap, nil
	}
	return nil, ErrJWTInvaild
}
