package Common

import (
	"crypto/rand"
	"errors"
	"math/big"
	"unicode"
)

const (
	charsetLower   = "abcdefghijklmnopqrstuvwxyz"
	charsetUpper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetDigits  = "0123456789"
	charsetSpecial = "!@#$%^&*()-_=+[]{}|;:'\",.<>?/`~"
)

var (
	ErrPasswordTooShort  = errors.New("رمز عبور باید حداقل ۸ کاراکتر باشد")
	ErrPasswordTooLong   = errors.New("رمز عبور نباید بیشتر از ۱۲۸ کاراکتر باشد")
	ErrNoUppercase       = errors.New("رمز عبور باید حداقل یک حرف بزرگ داشته باشد")
	ErrNoLowercase       = errors.New("رمز عبور باید حداقل یک حرف کوچک داشته باشد")
	ErrNoDigit           = errors.New("رمز عبور باید حداقل یک عدد داشته باشد")
	ErrNoSpecialChar     = errors.New("رمز عبور باید حداقل یک کاراکتر خاص داشته باشد (!@#$%^&* ...)")
	ErrPasswordHasSpaces = errors.New("رمز عبور نباید فاصله (space) داشته باشد")
	ErrPasswordIsCommon  = errors.New("این رمز عبور خیلی رایج و ناامن است")
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	if len(password) > 128 {
		return ErrPasswordTooLong
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
		hasSpace   bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsSpace(char):
			hasSpace = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if hasSpace {
		return ErrPasswordHasSpaces
	}

	if !hasUpper {
		return ErrNoUppercase
	}
	if !hasLower {
		return ErrNoLowercase
	}
	if !hasDigit {
		return ErrNoDigit
	}
	if !hasSpecial {
		return ErrNoSpecialChar
	}

	commonPasswords := map[string]bool{
		"12345678":   true,
		"password":   true,
		"123456789":  true,
		"qwertyuiop": true,
		"admin123":   true,
		"letmein":    true,
	}
	if commonPasswords[password] {
		return ErrPasswordIsCommon
	}

	return nil
}

func GeneratePassword(length int, includeSpecial bool) (string, error) {
	if length < 10 {
		return "", errors.New("طول رمز عبور باید حداقل ۱۰ باشد")
	}
	if length > 128 {
		return "", errors.New("طول رمز عبور نباید بیشتر از ۱۲۸ باشد")
	}

	charset := charsetLower + charsetUpper + charsetDigits
	if includeSpecial {
		charset += charsetSpecial
	}

	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}

	password := string(b)

	return password, nil
}
