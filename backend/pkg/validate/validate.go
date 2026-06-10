package validate

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	NameMax      = 50
	AddressMin   = 10
	AddressMax   = 300
	CityMax      = 100
	MaxAddresses = 20
)

const (
	phoneMinDigits = 10
	phoneMaxDigits = 15 // E.164
)

func phoneDigits(s string) string {
	var digits []rune
	for _, r := range s {
		if unicode.IsDigit(r) {
			digits = append(digits, r)
		}
	}
	return string(digits)
}

func Phone(s string, required bool) error {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		if required {
			return fmt.Errorf("phone is required")
		}
		return nil
	}
	digits := phoneDigits(trimmed)
	n := len(digits)
	if n < phoneMinDigits || n > phoneMaxDigits {
		return fmt.Errorf("invalid phone: expected %d..%d digits", phoneMinDigits, phoneMaxDigits)
	}
	return nil
}

func isNameRune(r rune) bool {
	return unicode.IsLetter(r) || r == ' ' || r == '-' || r == '\''
}

func PersonName(s string, required bool, fieldLabel string) error {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		if required {
			return fmt.Errorf("%s is required", fieldLabel)
		}
		return nil
	}
	if len([]rune(trimmed)) < 2 {
		return fmt.Errorf("%s is too short", fieldLabel)
	}
	if len([]rune(trimmed)) > NameMax {
		return fmt.Errorf("%s is too long", fieldLabel)
	}
	for _, r := range trimmed {
		if unicode.IsDigit(r) {
			return fmt.Errorf("%s must not contain digits", fieldLabel)
		}
		if !isNameRune(r) {
			return fmt.Errorf("%s contains invalid characters", fieldLabel)
		}
	}
	return nil
}

func FullName(s string) error {
	if err := PersonName(s, true, "full name"); err != nil {
		return err
	}
	parts := strings.Fields(strings.TrimSpace(s))
	if len(parts) < 2 {
		return fmt.Errorf("full name must contain at least two words")
	}
	return nil
}

func Address(s string, required bool) error {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		if required {
			return fmt.Errorf("address is required")
		}
		return nil
	}
	n := len([]rune(trimmed))
	if n < AddressMin {
		return fmt.Errorf("address is too short")
	}
	if n > AddressMax {
		return fmt.Errorf("address is too long")
	}
	return nil
}

func isCityRune(r rune) bool {
	return unicode.IsLetter(r) || r == ' ' || r == '-' || r == '.'
}

func City(s string) error {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return fmt.Errorf("city is required")
	}
	n := len([]rune(trimmed))
	if n < 2 {
		return fmt.Errorf("city is too short")
	}
	if n > CityMax {
		return fmt.Errorf("city is too long")
	}
	for _, r := range trimmed {
		if unicode.IsDigit(r) {
			return fmt.Errorf("city must not contain digits")
		}
		if !isCityRune(r) {
			return fmt.Errorf("city contains invalid characters")
		}
	}
	return nil
}

func UpdateUserProfile(firstName, lastName, phone *string) error {
	if firstName != nil && strings.TrimSpace(*firstName) != "" {
		if err := PersonName(*firstName, true, "first name"); err != nil {
			return err
		}
	}
	if lastName != nil && strings.TrimSpace(*lastName) != "" {
		if err := PersonName(*lastName, false, "last name"); err != nil {
			return err
		}
	}
	if phone != nil && strings.TrimSpace(*phone) != "" {
		if err := Phone(*phone, true); err != nil {
			return err
		}
	}
	return nil
}
