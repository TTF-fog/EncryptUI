package main

import (
	"math"
	"unicode"
)

// analyzePasswordStrength analyzes the strength of a password and returns a string descriptor.
// The strength is determined by calculating the password's entropy.
func analyzePasswordStrength(password string) string {
	if password == "" {
		return "very weak"
	}

	charSetSize := 0
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			if !hasLower {
				charSetSize += 26
				hasLower = true
			}
		case unicode.IsUpper(char):
			if !hasUpper {
				charSetSize += 26
				hasUpper = true
			}
		case unicode.IsDigit(char):
			if !hasDigit {
				charSetSize += 10
				hasDigit = true
			}
		default:
			// Using a common estimate for special characters.
			if !hasSpecial {
				charSetSize += 32
				hasSpecial = true
			}
		}
	}

	passwordLength := len(password)
	entropy := float64(passwordLength) * math.Log2(float64(charSetSize))

	switch {
	case entropy < 28:
		return "very weak"
	case entropy < 36:
		return "weak"
	case entropy < 60:
		return "moderate"
	case entropy < 128:
		return "strong"
	default:
		return "very strong"
	}
}
