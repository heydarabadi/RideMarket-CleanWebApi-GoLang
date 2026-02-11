package Common

import (
	"RideMarket-CleanWebApi-GoLang/Config"
	"crypto/rand"
	"errors"
	"fmt"
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

func GenerateOtp() string {
	cfg := Config.GetConfig()
	if cfg.Otp.Digits < 4 || cfg.Otp.Digits > 10 {
		panic("تعداد ارقام OTP باید بین ۴ تا ۱۰ باشد")
	}

	min := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(cfg.Otp.Digits-1)), nil)
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(cfg.Otp.Digits)), nil)
	max.Sub(max, big.NewInt(1))

	rangeSize := new(big.Int).Sub(max, min)
	rangeSize.Add(rangeSize, big.NewInt(1))

	n, err := rand.Int(rand.Reader, rangeSize)
	if err != nil {
		panic(err)
	}

	result := new(big.Int).Add(n, min)

	return fmt.Sprintf("%0*d", cfg.Otp.Digits, result)
}
