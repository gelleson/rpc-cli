package output

import (
	"testing"
)

func TestSensitiveMasker_IsSensitive(t *testing.T) {
	masker := NewSensitiveMasker()

	tests := []struct {
		name       string
		headerName string
		want       bool
	}{
		{
			name:       "authorization header",
			headerName: "Authorization",
			want:       true,
		},
		{
			name:       "bearer token",
			headerName: "Bearer-Token",
			want:       true,
		},
		{
			name:       "api key",
			headerName: "X-API-Key",
			want:       true,
		},
		{
			name:       "api key variant",
			headerName: "X-ApiKey",
			want:       true,
		},
		{
			name:       "password header",
			headerName: "X-Password",
			want:       true,
		},
		{
			name:       "secret header",
			headerName: "X-Secret",
			want:       true,
		},
		{
			name:       "auth token",
			headerName: "X-Auth-Token",
			want:       true,
		},
		{
			name:       "non-sensitive header",
			headerName: "Content-Type",
			want:       false,
		},
		{
			name:       "accept header",
			headerName: "Accept",
			want:       false,
		},
		{
			name:       "user agent",
			headerName: "User-Agent",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := masker.IsSensitive(tt.headerName)
			if got != tt.want {
				t.Errorf("IsSensitive(%s) = %v, want %v", tt.headerName, got, tt.want)
			}
		})
	}
}

func TestSensitiveMasker_Mask(t *testing.T) {
	masker := NewSensitiveMasker()

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "short value",
			value: "abc",
			want:  "****",
		},
		{
			name:  "exactly 4 chars",
			value: "abcd",
			want:  "****",
		},
		{
			name:  "long value",
			value: "Bearer_token_12345",
			want:  "Bear****",
		},
		{
			name:  "very long value",
			value: "this_is_a_very_long_secret_key_that_should_be_masked",
			want:  "this****",
		},
		{
			name:  "empty value",
			value: "",
			want:  "****",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := masker.Mask(tt.value)
			if got != tt.want {
				t.Errorf("Mask(%s) = %s, want %s", tt.value, got, tt.want)
			}
		})
	}
}

func TestSensitiveMasker_MaskIfSensitive(t *testing.T) {
	masker := NewSensitiveMasker()

	tests := []struct {
		name  string
		key   string
		value string
		want  string
	}{
		{
			name:  "sensitive key should be masked",
			key:   "Authorization",
			value: "Bearer token123456",
			want:  "Bear****",
		},
		{
			name:  "non-sensitive key should not be masked",
			key:   "Content-Type",
			value: "application/json",
			want:  "application/json",
		},
		{
			name:  "api key should be masked",
			key:   "X-API-Key",
			value: "secret_api_key_12345",
			want:  "secr****",
		},
		{
			name:  "accept header should not be masked",
			key:   "Accept",
			value: "application/json",
			want:  "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := masker.MaskIfSensitive(tt.key, tt.value)
			if got != tt.want {
				t.Errorf("MaskIfSensitive(%s, %s) = %s, want %s", tt.key, tt.value, got, tt.want)
			}
		})
	}
}

func TestSensitiveMasker_CaseInsensitive(t *testing.T) {
	masker := NewSensitiveMasker()

	// Test that matching is case-insensitive
	headers := []string{
		"authorization",
		"Authorization",
		"AUTHORIZATION",
		"AuThOrIzAtIoN",
	}

	for _, header := range headers {
		if !masker.IsSensitive(header) {
			t.Errorf("IsSensitive should be case-insensitive, failed for: %s", header)
		}
	}
}
