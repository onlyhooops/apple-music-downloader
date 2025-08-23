package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"main/internal/core"
)

// EnsureSafePath truncates path components to ensure the total path length does not exceed the limit
func EnsureSafePath(basePath, artistDir, albumDir, fileName string) (string, string, string) {
	truncate := func(s string, n int) string {
		if n <= 0 {
			return s
		}
		runes := []rune(s)
		if len(runes) <= n {
			return ""
		}
		return string(runes[:len(runes)-n])
	}

	for {
		currentPath := filepath.Join(basePath, artistDir, albumDir, fileName)
		if len(currentPath) <= core.MaxPathLength {
			break
		}

		overage := len(currentPath) - core.MaxPathLength
		ext := filepath.Ext(fileName)
		stem := strings.TrimSuffix(fileName, ext)

		var prefixPart string
		var namePart string
		re := regexp.MustCompile(`^(\d+[\s.-]*)`)
		matches := re.FindStringSubmatch(stem)

		if len(matches) > 1 {
			prefixPart = matches[1]
			namePart = strings.TrimPrefix(stem, prefixPart)
		} else {
			prefixPart = ""
			namePart = stem
		}

		if len(namePart) > 0 {
			canShorten := len(namePart)
			shortenAmount := overage
			if shortenAmount > canShorten {
				shortenAmount = canShorten
			}
			namePart = truncate(namePart, shortenAmount)

			if namePart == "" {
				prefixPart = strings.TrimRight(prefixPart, " .-")
			}

			fileName = prefixPart + namePart + ext
			continue
		}

		if len(albumDir) > 1 { // 至少保留一个字符
			canShorten := len(albumDir)
			shortenAmount := overage
			if shortenAmount > canShorten {
				shortenAmount = canShorten
			}
			albumDir = truncate(albumDir, shortenAmount)
			continue
		}

		if len(artistDir) > 1 { // 至少保留一个字符
			canShorten := len(artistDir)
			shortenAmount := overage
			if shortenAmount > canShorten {
				shortenAmount = canShorten
			}
			artistDir = truncate(artistDir, shortenAmount)
			continue
		}

		break
	}

	return artistDir, albumDir, fileName
}

// IsInArray checks if a target integer is in an array of integers
func IsInArray(arr []int, target int) bool {
	for _, num := range arr {
		if num == target {
			return true
		}
	}
	return false
}

// FileExists checks if a file exists at the given path
func FileExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		return !f.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FormatSpeed formats a download speed from bytes/sec to a human-readable string
func FormatSpeed(bytesPerSecond float64) string {
	if bytesPerSecond < 1024 {
		return fmt.Sprintf("%.1f B/s", bytesPerSecond)
	}
	kbps := bytesPerSecond / 1024
	if kbps < 1024 {
		return fmt.Sprintf("%.1f KB/s", kbps)
	}
	mbps := kbps / 1024
	return fmt.Sprintf("%.1f MB/s", mbps)
}

// Contains checks if a string slice contains a specific item
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

