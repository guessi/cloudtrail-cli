package utils

import (
	"testing"

	"github.com/guessi/cloudtrail-cli/pkg/types"
)

func TestIsValidEventSource(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid AWS service domain",
			input:    "s3.amazonaws.com",
			expected: true,
		},
		{
			name:     "Invalid service name",
			input:    "invalid",
			expected: false,
		},
		{
			name:     "Empty string input",
			input:    "",
			expected: false,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := isValidEventSource(tc.input)
			if got != tc.expected {
				t.Errorf("isValidEventSource(%q) = %v, want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestIsValidUUID(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid UUID format",
			input:    "550e8400-e29b-41d4-a716-446655440000",
			expected: true,
		},
		{
			name:     "Invalid UUID format",
			input:    "invalid-uuid",
			expected: false,
		},
		{
			name:     "Empty string input",
			input:    "",
			expected: false,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			got := isValidUUID(tc.input)
			if got != tc.expected {
				t.Errorf("isValidUUID(%q) = %v, want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestIsValidAccessKeyID(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid IAM access key (AKIA prefix)",
			input:    "AKIAIOSFODNN7EXAMPLE",
			expected: true,
		},
		{
			name:     "Valid STS access key (ASIA prefix)",
			input:    "ASIAIOSFODNN7EXAMPLE",
			expected: true,
		},
		{
			name:     "Invalid access key format",
			input:    "invalid",
			expected: false,
		},
		{
			name:     "Empty string input",
			input:    "",
			expected: false,
		},
		{
			name:     "Too short access key",
			input:    "AKIA123",
			expected: false,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			got := isValidAccessKeyID(tc.input)
			if got != tc.expected {
				t.Errorf("isValidAccessKeyID(%q) = %v, want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	testCases := []struct {
		name     string
		truncate bool
		input    string
		length   int
		expected string
	}{
		{
			name:     "Truncation disabled",
			truncate: false,
			input:    "hello world",
			length:   5,
			expected: "hello world",
		},
		{
			name:     "Short string",
			truncate: false,
			input:    "hello",
			length:   10,
			expected: "hello",
		},
		{
			name:     "Long string truncated",
			truncate: true,
			input:    "hello world",
			length:   5,
			expected: "hello",
		},
		{
			name:     "Exact length",
			truncate: true,
			input:    "hello",
			length:   5,
			expected: "hello",
		},
		{
			name:     "Empty string",
			truncate: true,
			input:    "",
			length:   5,
			expected: "",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			got := truncateString(tc.truncate, tc.input, tc.length)
			if got != tc.expected {
				t.Errorf("truncateString(%v, %q, %d) = %q, want %q", tc.truncate, tc.input, tc.length, got, tc.expected)
			}
		})
	}
}

func TestGetDisplayUserName(t *testing.T) {
	testCases := []struct {
		name     string
		identity types.UserIdentity
		expected string
	}{
		{
			name: "IAM User",
			identity: types.UserIdentity{
				Type:     "IAMUser",
				UserName: "testuser",
			},
			expected: "testuser",
		},
		{
			name: "Assumed Role",
			identity: types.UserIdentity{
				Type: "AssumedRole",
				Arn:  "arn:aws:sts::123456789012:assumed-role/MyRole/session",
			},
			expected: "session",
		},
		{
			name: "Root User",
			identity: types.UserIdentity{
				Type:      "Root",
				AccountId: "123456789012",
			},
			expected: "-",
		},
		{
			name: "Unknown Type",
			identity: types.UserIdentity{
				Type: "Unknown",
			},
			expected: "-",
		},
		{
			name: "WebIdentity User",
			identity: types.UserIdentity{
				Type:     "WebIdentityUser",
				UserName: "webuser",
			},
			expected: "webuser",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			got := getDisplayUserName(tc.identity)
			if got != tc.expected {
				t.Errorf("getDisplayUserName(%+v) = %q, want %q", tc.identity, got, tc.expected)
			}
		})
	}
}

func TestGetBatchSize(t *testing.T) {
	testCases := []struct {
		name     string
		input    int
		expected *int32
	}{
		{
			name:     "Valid size",
			input:    25,
			expected: func() *int32 { v := int32(25); return &v }(),
		},
		{
			name:     "Zero size uses default",
			input:    0,
			expected: func() *int32 { v := int32(50); return &v }(),
		},
		{
			name:     "Negative size uses default",
			input:    -1,
			expected: func() *int32 { v := int32(50); return &v }(),
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			got := getBatchSize(tc.input)
			if got == nil || tc.expected == nil {
				if got != tc.expected {
					t.Errorf("getBatchSize(%d) = %v, want %v", tc.input, got, tc.expected)
				}
			} else if *got != *tc.expected {
				t.Errorf("getBatchSize(%d) = %d, want %d", tc.input, *got, *tc.expected)
			}
		})
	}
}

func TestBuildLookupAttributes(t *testing.T) {
	testCases := []struct {
		name     string
		input    types.CloudTrailCliInput
		expected int
	}{
		{
			name:     "No filters",
			input:    types.CloudTrailCliInput{},
			expected: 0,
		},
		{
			name: "Single filter",
			input: types.CloudTrailCliInput{
				EventName: "GetObject",
			},
			expected: 1,
		},
		{
			name: "Multiple filters",
			input: types.CloudTrailCliInput{
				EventName: "GetObject",
				UserName:  "testuser",
			},
			expected: 2,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			got := buildLookupAttributes(tc.input)
			if len(got) != tc.expected {
				t.Errorf("buildLookupAttributes() returned %d filters, want %d", len(got), tc.expected)
			}
		})
	}
}
