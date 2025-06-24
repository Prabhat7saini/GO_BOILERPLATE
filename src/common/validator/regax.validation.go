package validator

import (
	"regexp"
	"github.com/go-playground/validator/v10"
)

var (
	uppercaseRegex   = regexp.MustCompile(`[A-Z]`)
	lowercaseRegex   = regexp.MustCompile(`[a-z]`)
	digitRegex       = regexp.MustCompile(`[0-9]`)
	specialCharRegex = regexp.MustCompile(`[\W_]`)
)

func StrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	length := len(password)

	return length >= 8 && length <= 12 &&
		uppercaseRegex.MatchString(password) &&
		lowercaseRegex.MatchString(password) &&
		digitRegex.MatchString(password) &&
		specialCharRegex.MatchString(password)
}



// func StrongPassword(fl validator.FieldLevel) bool {
// 	password := fl.Field().String()

// 	// At least one uppercase, one lowercase, one digit, one special char, and 8+ chars
// 	regex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z\d]).{8,}$`
// 	match, _ := regexp.MatchString(regex, password)

// 	return match
// }