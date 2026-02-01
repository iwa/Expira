package state

import (
	"testing"
	"time"
)

func TestDomainGetDaysUntilExpiryShouldReturnCorrectValue(t *testing.T) {
	domain := Domain{
		Name:       "example.org",
		Exists:     true,
		ExpiryDate: time.Now().Add(30 * 24 * time.Hour), // 30 days after now
	}

	result := domain.GetDaysUntilExpiry()

	if result != 30 {
		t.Fatal("wrong result, got", result, "instead of 30")
	}
}

func TestDomainGetDaysUntilExpiryShouldReturnZeroOnExpiryDay(t *testing.T) {
	domain := Domain{
		Name:       "example.org",
		Exists:     true,
		ExpiryDate: time.Now(),
	}

	result := domain.GetDaysUntilExpiry()

	if result != 0 {
		t.Fatal("wrong result, got", result, "instead of 0")
	}
}

func TestDomainGetDaysUntilExpiryShouldReturnNegativeForExpired(t *testing.T) {
	domain := Domain{
		Name:       "example.org",
		Exists:     true,
		ExpiryDate: time.Now().Add(-1 * 30 * 24 * time.Hour), // 30 days before now
	}

	result := domain.GetDaysUntilExpiry()

	if result >= 0 {
		t.Fatal("wrong result, got", result, "instead of negative value")
	}
}
