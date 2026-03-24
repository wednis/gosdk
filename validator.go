package gosdk

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	passwordSpecialSet = func() map[rune]bool {
		set := make(map[rune]bool)
		for _, r := range ",.!?;:_-+=*/#%$@&^|()[]{}<>~`" {
			set[r] = true
		}
		return set
	}()

	emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// 验证用户名是否合法
//   - unicode字符数 3 - 16（含）
//   - 特殊字符仅允许 '-' '_'
//   - 允许大小写英文字母
//   - 允许简繁中文字符
//   - 允许数字
//   - 无首字符限制
func ValidateUsername(username string) bool {
	count := utf8.RuneCountInString(username)
	if count < 3 || count > 16 {
		return false
	}

	for _, r := range username {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
		case r >= '0' && r <= '9':
		case r == '-' || r == '_':
		case unicode.Is(unicode.Han, r):
		default:
			return false
		}
	}
	return true
}

// 验证密码是否合法
//   - 长度 8 - 20（含）
//   - 需要包含小写英文字母，大写英文字母，数字，特殊字符
//   - 特殊字符串为 "@!_-#$*.?&%=+/:;[]{}()<>|^`~"  出于安全考虑没有单双引号
func ValidatePassword(password string) bool {
	count := utf8.RuneCountInString(password)
	if count < 8 || count > 20 {
		return false
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false
	for _, r := range password {
		switch {
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= '0' && r <= '9':
			hasDigit = true
		case passwordSpecialSet[r]:
			hasSpecial = true
		default:
			// 出现非法字符
			return false
		}
	}
	return hasLower && hasUpper && hasDigit && hasSpecial
}

// 验证邮箱是否合法
//   - 长度 5 - 254（含）
//   - 正则匹配规则 `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
func ValidateEmail(email string) bool {
	count := len(email)
	if count < 5 || count > 254 {
		return false
	}

	// 快速检查
	// 必须包含 @，且不能以@开头或结尾
	if !strings.Contains(email, "@") || strings.HasPrefix(email, "@") || strings.HasSuffix(email, "@") {
		return false
	}
	return emailRegexp.MatchString(email)
}
