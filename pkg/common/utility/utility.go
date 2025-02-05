package utility

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"regexp"
	"strings"
	"unicode"
)

func IsValidPhoneNumber(phone string) bool {
	// Simple regex pattern for basic phone number validation
	fmt.Println(" check pfone validity")
	const phoneRegex = `^\+?[1-9]\d{1,14}$` // E.164 international phone number format
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}

func ParseAndValidate(c *fiber.Ctx, request interface{}, validate *validator.Validate) (err error) {
	if err = c.BodyParser(request); err != nil {
		return
	}
	if err = validate.Struct(request); err != nil {
		return
	}
	return nil
}
func IsValidPassword(password string) bool {

	if len(password) < 9 {
		return false
	}

	specialCharPattern := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	if !specialCharPattern.MatchString(password) {
		return false
	}
	var hasLetter, hasDigit, hasUppercase bool
	for _, char := range password {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
			if unicode.IsUpper(char) {
				hasUppercase = true
			}
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit || !hasUppercase {
		return false
	}
	if strings.Contains(password, " ") {
		return false
	}

	return true
}
