package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type FileFormat int

const (
	FormatYAML FileFormat = iota
	FormatJSON
	FormatUnknown
)

func DetectFormat(filename string) FileFormat {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".yaml", ".yml":
		return FormatYAML
	case ".json":
		return FormatJSON
	default:
		return FormatUnknown
	}
}

func DetectFormatFromContent(data []byte) FileFormat {
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 {
		return FormatUnknown
	}
	if trimmed[0] == '{' || trimmed[0] == '[' {
		return FormatJSON
	}
	return FormatYAML
}

func Unmarshal(data []byte, format FileFormat, v interface{}) error {
	if len(data) == 0 {
		return nil
	}
	switch format {
	case FormatJSON:
		return json.Unmarshal(data, v)
	case FormatYAML:
		return yaml.Unmarshal(data, v)
	default:
		return fmt.Errorf("unknown file format")
	}
}

func Marshal(v interface{}, format FileFormat) ([]byte, error) {
	switch format {
	case FormatJSON:
		return json.MarshalIndent(v, "", "  ")
	case FormatYAML:
		return yaml.Marshal(v)
	default:
		return nil, fmt.Errorf("unknown file format")
	}
}

func GetAlternateFilename(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	switch strings.ToLower(ext) {
	case ".yaml", ".yml":
		return base + ".json"
	case ".json":
		return base + ".yaml"
	default:
		return ""
	}
}
