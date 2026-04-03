package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBootEntry_MatchesOS(t *testing.T) {
	tests := []struct {
		name     string
		entry    BootEntry
		osName   OSName
		expected bool
	}{
		{
			name:     "exact match",
			entry:    BootEntry{Name: "Ubuntu"},
			osName:   "Ubuntu",
			expected: true,
		},
		{
			name:     "case insensitive match",
			entry:    BootEntry{Name: "ubuntu (recovery mode)"},
			osName:   "Ubuntu",
			expected: true,
		},
		{
			name:     "substring match",
			entry:    BootEntry{Name: "Windows 10"},
			osName:   "Windows",
			expected: true,
		},
		{
			name:     "no match",
			entry:    BootEntry{Name: "Ubuntu"},
			osName:   "Windows",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.entry.MatchesOS(tt.osName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMatchGrubEntryToOS(t *testing.T) {
	entries := []BootEntry{
		{Name: "Ubuntu"},
		{Name: "Windows 10"},
		{Name: "Windows 11"},
		{Name: "Advanced options for Ubuntu"},
	}

	tests := []struct {
		name          string
		osName        OSName
		expectedName  string
		expectedError bool
	}{
		{
			name:         "find Ubuntu",
			osName:       "Ubuntu",
			expectedName: "Ubuntu",
		},
		{
			name:         "find Windows (first match)",
			osName:       "Windows",
			expectedName: "Windows 10",
		},
		{
			name:          "not found",
			osName:        "MacOS",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MatchGrubEntryToOS(entries, tt.osName)
			if tt.expectedError {
				assert.Error(t, err)
				assert.True(t, strings.Contains(err.Error(), string(tt.osName)), "error should contain OS name")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedName, result)
			}
		})
	}
}
