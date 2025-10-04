package output

import "strings"

// SensitiveMasker handles masking of sensitive information
type SensitiveMasker struct {
	sensitiveKeywords []string
}

// NewSensitiveMasker creates a new SensitiveMasker
func NewSensitiveMasker() *SensitiveMasker {
	return &SensitiveMasker{
		sensitiveKeywords: []string{
			"authorization",
			"token",
			"api-key",
			"apikey",
			"x-api-key",
			"x-auth-token",
			"secret",
			"password",
			"bearer",
		},
	}
}

// IsSensitive checks if a header name is sensitive
func (m *SensitiveMasker) IsSensitive(headerName string) bool {
	lowerName := strings.ToLower(headerName)
	for _, keyword := range m.sensitiveKeywords {
		if strings.Contains(lowerName, keyword) {
			return true
		}
	}
	return false
}

// Mask masks a sensitive value
func (m *SensitiveMasker) Mask(value string) string {
	if len(value) <= 4 {
		return "****"
	}
	return value[:4] + "****"
}

// MaskIfSensitive masks a value if the key is sensitive
func (m *SensitiveMasker) MaskIfSensitive(key, value string) string {
	if m.IsSensitive(key) {
		return m.Mask(value)
	}
	return value
}
