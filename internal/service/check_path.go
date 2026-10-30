package service

import (
	"net/url"
	"path/filepath"
	"strings"
	"unicode"
)

func IsURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	return err == nil &&
		(u.Scheme == "http" || u.Scheme == "https") &&
		u.Host != ""
}

func IsFileName(s string) bool {
	if s == "" || strings.TrimSpace(s) != s {
		return false
	}

	if s == "." || s == ".." {
		return false
	}

	if strings.ContainsAny(s, `<>:"|?*`) {
		return false
	}

	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}

	return filepath.Base(s) == s
}
