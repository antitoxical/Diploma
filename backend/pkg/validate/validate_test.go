package validate_test

import (
	"testing"

	"telegramshop_backend/pkg/validate"
)

func TestPhone(t *testing.T) {
	if err := validate.Phone("+375291234567", true); err != nil {
		t.Fatalf("expected valid phone, got %v", err)
	}
	if err := validate.Phone("123", true); err == nil {
		t.Fatal("expected invalid phone")
	}
}

func TestPersonNameRejectsDigits(t *testing.T) {
	if err := validate.PersonName("Ivan123", true, "name"); err == nil {
		t.Fatal("expected digits to be rejected")
	}
}
