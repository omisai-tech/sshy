package cmd

import "testing"

func TestHasAllTags(t *testing.T) {
	tests := []struct {
		name       string
		serverTags []string
		filterTags []string
		expected   bool
	}{
		{
			name:       "empty filter tags",
			serverTags: []string{"prod", "web"},
			filterTags: []string{},
			expected:   true,
		},
		{
			name:       "single matching tag",
			serverTags: []string{"prod", "web"},
			filterTags: []string{"prod"},
			expected:   true,
		},
		{
			name:       "multiple matching tags",
			serverTags: []string{"prod", "web", "api"},
			filterTags: []string{"prod", "web"},
			expected:   true,
		},
		{
			name:       "all tags match",
			serverTags: []string{"prod", "web"},
			filterTags: []string{"prod", "web"},
			expected:   true,
		},
		{
			name:       "missing tag",
			serverTags: []string{"prod", "web"},
			filterTags: []string{"staging"},
			expected:   false,
		},
		{
			name:       "partial match",
			serverTags: []string{"prod", "web"},
			filterTags: []string{"prod", "staging"},
			expected:   false,
		},
		{
			name:       "empty server tags",
			serverTags: []string{},
			filterTags: []string{"prod"},
			expected:   false,
		},
		{
			name:       "both empty",
			serverTags: []string{},
			filterTags: []string{},
			expected:   true,
		},
		{
			name:       "nil server tags",
			serverTags: nil,
			filterTags: []string{"prod"},
			expected:   false,
		},
		{
			name:       "nil filter tags",
			serverTags: []string{"prod"},
			filterTags: nil,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasAllTags(tt.serverTags, tt.filterTags)
			if result != tt.expected {
				t.Errorf("hasAllTags(%v, %v) = %v, expected %v", tt.serverTags, tt.filterTags, result, tt.expected)
			}
		})
	}
}
