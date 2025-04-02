package utils

import (
	"net/http"
	"regexp"
	"strings"
)

func IsValidEmail(email string) bool {
	// Regular expression for validating an email
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
func IsValidPassword(password string) error {
	minLength := 8
	maxLength := 128
	if len(password) < minLength || len(password) > maxLength {
		return NewUserError(http.StatusBadRequest, "INVALID_LENGTH")
	}

	// Проверка на наличие пробельных символов
	if strings.ContainsAny(password, " \t\n\r") {
		return NewUserError(http.StatusBadRequest, "SPACES_NOT_ALLOWED")
	}

	// Проверка на наличие как минимум одной цифры
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return NewUserError(http.StatusBadRequest, "DIGIT_REQUIRED")
	}

	// Проверка на наличие как минимум одной латинской буквы в нижнем регистре
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return NewUserError(http.StatusBadRequest, "LOWERCASE_REQUIRED")
	}

	// Проверка на наличие как минимум одной латинской буквы в верхнем регистре
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return NewUserError(http.StatusBadRequest, "UPPERCASE_REQUIRED")
	}

	// Проверка на наличие как минимум одного специального символа
	specialChars := regexp.QuoteMeta("-!@#$%^&*()_=+[]{};:',.<>?")
	if !regexp.MustCompile(`[` + specialChars + `]`).MatchString(password) {
		return NewUserError(http.StatusBadRequest, "SPECIAL_SYMBOL_REQUIRED")
	}

	return nil
}
