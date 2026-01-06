package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeSlug creates a URL-friendly slug from a string
// Rules:
// - Convert to lowercase
// - Remove accents and special characters
// - Replace spaces with hyphens
// - Remove consecutive hyphens
// - Trim hyphens from start and end
// Example: "Workshop Bronkoskopi Tingkat Lanjut" -> "workshop-bronkoskopi-tingkat-lanjut"
func NormalizeSlug(input string) string {
	if input == "" {
		return ""
	}

	// Step 1: Convert to lowercase
	slug := strings.ToLower(input)

	// Step 2: Remove accents and normalize unicode
	slug = removeAccents(slug)

	// Step 3: Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Step 4: Remove any character that's not alphanumeric or hyphen
	reg := regexp.MustCompile("[^a-z0-9-]")
	slug = reg.ReplaceAllString(slug, "")

	// Step 5: Replace consecutive hyphens with single hyphen
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Step 6: Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// removeAccents removes diacritical marks from unicode characters
// Example: "café" -> "cafe", "naïve" -> "naive"
func removeAccents(input string) string {
	// NFD (Canonical Decomposition) separates base characters from accents
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn = Nonspacing_Mark (accents)
	}))

	result, _, _ := transform.String(t, input)
	return result
}

// IsSlugUnique checks if a slug is unique in a given table
// This is a helper function - actual uniqueness check should be done in the model/database layer
func IsSlugUnique(slug string) bool {
	return slug != "" && len(slug) > 0
}

// GenerateUniqueSlug generates a unique slug by appending a number if needed
// Usage: When a slug already exists, append "-2", "-3", etc.
// Example: "workshop-bronkoskopi", "workshop-bronkoskopi-2", "workshop-bronkoskopi-3"
func GenerateUniqueSlug(baseSlug string, existingSlugs map[string]bool) string {
	if !existingSlugs[baseSlug] {
		return baseSlug
	}

	// Try with numbers: baseSlug-2, baseSlug-3, etc.
	for i := 2; i <= 1000; i++ {
		candidateSlug := baseSlug + "-" + fmt.Sprintf("%d", i)
		if !existingSlugs[candidateSlug] {
			return candidateSlug
		}
	}

	// Fallback: append timestamp-based suffix if too many conflicts
	return baseSlug + "-" + fmt.Sprintf("%.0f", GetCurrentUnixTimestamp())
}

// ValidateSlug checks if a slug is valid
// A valid slug should:
// - Not be empty
// - Be lowercase
// - Only contain alphanumeric characters and hyphens
// - Not start or end with hyphen
func ValidateSlug(slug string) bool {
	if slug == "" {
		return false
	}

	// Check if it starts or ends with hyphen
	if slug[0] == '-' || slug[len(slug)-1] == '-' {
		return false
	}

	// Check if it contains only valid characters
	validSlugRegex := regexp.MustCompile("^[a-z0-9-]+$")
	return validSlugRegex.MatchString(slug)
}

// GenerateSlug is an alias for NormalizeSlug for better readability
// Creates a URL-friendly slug from a string
func GenerateSlug(input string) string {
	return NormalizeSlug(input)
}
