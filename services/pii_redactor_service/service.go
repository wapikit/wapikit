package pii_redactor

import (
	"regexp"

	"github.com/cockroachdb/redact"
)

type PIIRedactor struct {
	patterns []*regexp.Regexp
}

func NewPIIRedactor() *PIIRedactor {
	// Define regex patterns for common PII types
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),                               // SSN
		regexp.MustCompile(`\b\d{4} \d{4} \d{4} \d{4}\b`),                         // Credit Card
		regexp.MustCompile(`\b\d{10}\b`),                                          // Phone Number
		regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`), // Email
	}
	return &PIIRedactor{patterns: patterns}
}

func (r *PIIRedactor) Redact(input string) string {
	redacted := redact.RedactableString(input)
	for _, pattern := range r.patterns {
		redacted = redact.RedactableString(pattern.ReplaceAllString(string(redacted), "[REDACTED]"))
	}
	return string(redacted.Redact())
}
