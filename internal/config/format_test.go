package config

import (
	"testing"
)

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		filename string
		expected FileFormat
	}{
		{"servers.yaml", FormatYAML},
		{"servers.yml", FormatYAML},
		{"servers.YAML", FormatYAML},
		{"servers.json", FormatJSON},
		{"servers.JSON", FormatJSON},
		{"servers.txt", FormatUnknown},
		{"servers", FormatUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := DetectFormat(tt.filename)
			if result != tt.expected {
				t.Errorf("DetectFormat(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestDetectFormatFromContent(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected FileFormat
	}{
		{"JSON object", []byte(`{"name": "test"}`), FormatJSON},
		{"JSON array", []byte(`[{"name": "test"}]`), FormatJSON},
		{"JSON with whitespace", []byte(`  {"name": "test"}`), FormatJSON},
		{"YAML", []byte(`name: test`), FormatYAML},
		{"YAML list", []byte(`- name: test`), FormatYAML},
		{"Empty", []byte{}, FormatUnknown},
		{"Whitespace only", []byte("   "), FormatUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectFormatFromContent(tt.content)
			if result != tt.expected {
				t.Errorf("DetectFormatFromContent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	type TestStruct struct {
		Name string `yaml:"name" json:"name"`
		Port int    `yaml:"port" json:"port"`
	}

	original := TestStruct{Name: "test", Port: 22}

	t.Run("YAML roundtrip", func(t *testing.T) {
		data, err := Marshal(original, FormatYAML)
		if err != nil {
			t.Fatalf("Marshal YAML error: %v", err)
		}

		var result TestStruct
		err = Unmarshal(data, FormatYAML, &result)
		if err != nil {
			t.Fatalf("Unmarshal YAML error: %v", err)
		}

		if result.Name != original.Name || result.Port != original.Port {
			t.Errorf("YAML roundtrip failed: got %+v, want %+v", result, original)
		}
	})

	t.Run("JSON roundtrip", func(t *testing.T) {
		data, err := Marshal(original, FormatJSON)
		if err != nil {
			t.Fatalf("Marshal JSON error: %v", err)
		}

		var result TestStruct
		err = Unmarshal(data, FormatJSON, &result)
		if err != nil {
			t.Fatalf("Unmarshal JSON error: %v", err)
		}

		if result.Name != original.Name || result.Port != original.Port {
			t.Errorf("JSON roundtrip failed: got %+v, want %+v", result, original)
		}
	})

	t.Run("Unknown format marshal error", func(t *testing.T) {
		_, err := Marshal(original, FormatUnknown)
		if err == nil {
			t.Error("Expected error for unknown format marshal")
		}
	})

	t.Run("Unknown format unmarshal error", func(t *testing.T) {
		var result TestStruct
		err := Unmarshal([]byte(`{"name":"test"}`), FormatUnknown, &result)
		if err == nil {
			t.Error("Expected error for unknown format unmarshal")
		}
	})
}

func TestGetAlternateFilename(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"servers.yaml", "servers.json"},
		{"servers.yml", "servers.json"},
		{"servers.json", "servers.yaml"},
		{"config.yaml", "config.json"},
		{"servers.txt", ""},
		{"servers", ""},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := GetAlternateFilename(tt.filename)
			if result != tt.expected {
				t.Errorf("GetAlternateFilename(%q) = %q, want %q", tt.filename, result, tt.expected)
			}
		})
	}
}
