package generictypes_test

import (
	"testing"

	"github.com/giantswarm/generic-types-go"
)

func TestDomainValidatorValidDomain(t *testing.T) {
	d := generictypes.Domain("i.am.correct.com")

	if err := d.Validate(); err != nil {
		t.Fatalf("Valid domain detected to be invalid: %v", err)
	}
}

func TestDomainValidatorInvalidDomain(t *testing.T) {
	d := generictypes.Domain("i.$am.invalid.com")

	if err := d.Validate(); err == nil {
		t.Fatalf("Invalid domain detected to be valid: %v", d.String())
	}
}

func TestDomainValidatorInvalidDomainWithPort(t *testing.T) {
	d := generictypes.Domain("i.am.invalid.com:80")

	if err := d.Validate(); err == nil {
		t.Fatalf("Invalid domain detected to be valid: %v", d.String())
	}
}
