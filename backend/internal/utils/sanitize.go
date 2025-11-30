package utils

import (
	"html"
	"strings"
)

// SanitizeString sanitizes user input to prevent XSS
func SanitizeString(s string) string {
	// Escape HTML entities
	s = html.EscapeString(s)
	// Trim whitespace
	s = strings.TrimSpace(s)
	return s
}

// SanitizeStrings sanitizes array of strings
func SanitizeStrings(arr []string) []string {
	if arr == nil {
		return []string{}
	}
	
	result := make([]string, 0, len(arr))
	for _, s := range arr {
		sanitized := SanitizeString(s)
		// Skip empty strings after sanitization
		if sanitized != "" {
			result = append(result, sanitized)
		}
	}
	return result
}

// SanitizeURL basic URL sanitization (allows only http/https)
func SanitizeURL(url string) string {
	url = strings.TrimSpace(url)
	
	// Empty URL is OK
	if url == "" {
		return ""
	}
	
	// Must start with http:// or https://
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return ""
	}
	
	// Basic check for common XSS patterns
	lowerURL := strings.ToLower(url)
	if strings.Contains(lowerURL, "javascript:") || 
	   strings.Contains(lowerURL, "data:") || 
	   strings.Contains(lowerURL, "vbscript:") {
		return ""
	}
	
	return url
}

// TruncateString truncates string to max length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}