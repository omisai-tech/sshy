package config

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/omisai-tech/sshy/internal/models"
)

const (
	DefaultHTTPTimeout = 30 * time.Second
)

var httpClient = &http.Client{
	Timeout: DefaultHTTPTimeout,
}

func IsURL(path string) bool {
	if path == "" {
		return false
	}
	lower := strings.ToLower(path)
	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}

func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}
	if parsed.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	return nil
}

func FetchServersFromURL(urlStr string) (models.Servers, error) {
	if err := ValidateURL(urlStr); err != nil {
		return nil, err
	}

	resp, err := httpClient.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if len(data) == 0 {
		return models.Servers{}, nil
	}

	format := DetectFormatFromContent(data)
	if format == FormatUnknown {
		format = detectFormatFromContentType(resp.Header.Get("Content-Type"))
	}
	if format == FormatUnknown {
		format = detectFormatFromURL(urlStr)
	}
	if format == FormatUnknown {
		format = FormatYAML
	}

	var servers models.Servers
	if err := Unmarshal(data, format, &servers); err != nil {
		return nil, fmt.Errorf("failed to parse servers data: %w", err)
	}

	return servers, nil
}

func FetchServersFromURLWithSource(urlStr string) ([]models.ServerWithSource, error) {
	servers, err := FetchServersFromURL(urlStr)
	if err != nil {
		return nil, err
	}

	result := make([]models.ServerWithSource, len(servers))
	for i, server := range servers {
		result[i] = models.ServerWithSource{
			Server: server,
			Source: models.SourceShared,
		}
	}

	return result, nil
}

func detectFormatFromContentType(contentType string) FileFormat {
	lower := strings.ToLower(contentType)
	if strings.Contains(lower, "json") {
		return FormatJSON
	}
	if strings.Contains(lower, "yaml") || strings.Contains(lower, "yml") {
		return FormatYAML
	}
	return FormatUnknown
}

func detectFormatFromURL(urlStr string) FileFormat {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return FormatUnknown
	}
	path := strings.ToLower(parsed.Path)
	if strings.HasSuffix(path, ".json") {
		return FormatJSON
	}
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return FormatYAML
	}
	return FormatUnknown
}
