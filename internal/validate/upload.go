package validate

import (
	"errors"
	"mime"
	"path/filepath"
	"strings"
)

const MaxBytes = 20 * 1024 * 1024

var types = map[string]string{
	".pdf": "application/pdf",
}

func CheckFile(name, contentType string, size int64) error {
	if size <= 0 {
		return errors.New("File is empty")
	}
	if size > MaxBytes {
		return errors.New("File is larger than 20 MB")
	}
	ext := strings.ToLower(filepath.Ext(filepath.Base(name)))
	want, ok := types[ext]
	if !ok {
		return errors.New("Only PDF files are allowed")
	}
	media := strings.TrimSpace(contentType)
	if media == "" || strings.EqualFold(media, "application/octet-stream") {
		return nil
	}
	parsed, _, err := mime.ParseMediaType(media)
	if err != nil || !strings.EqualFold(parsed, want) {
		return errors.New("File type does not match its extension")
	}
	return nil
}

var labels = map[string]bool{
	"petition": true, "order": true, "evidence": true, "vakalatnama": true, "other": true,
}

func Label(label string) bool { return labels[label] }

var caseTypes = map[string]bool{
	"civil": true, "criminal": true, "writ": true, "family": true, "consumer": true, "other": true,
}

func CaseType(v string) bool { return caseTypes[v] }

var caseStatuses = map[string]bool{"active": true, "stayed": true, "disposed": true}

func CaseStatus(v string) bool { return caseStatuses[v] }

var purposes = map[string]bool{
	"mention": true, "arguments": true, "evidence": true, "orders": true, "compliance": true, "other": true,
}

func Purpose(v string) bool { return purposes[v] }

func Email(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" {
		return true
	}
	at := strings.LastIndex(v, "@")
	if at <= 0 || at != strings.Index(v, "@") || strings.ContainsAny(v, " \t") {
		return false
	}
	return strings.Contains(v[at+1:], ".")
}
