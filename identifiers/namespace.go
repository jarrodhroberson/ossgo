package identifiers

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const PATH_SEPARATOR = "/"
const LABEL_SEPARATOR = ":"
const NAMESPACE_LABEL_SEPARATOR = "?"
const KEY_VALUE_SEPARATOR = "="

type Path string

func (p Path) String() string {
	return string(p)
}

type Entry struct {
	Path       string            `json:"path"`
	Labels     map[string]string `json:"labels"`
	Data       []byte            `json:"data"`
	Expiration time.Duration     `json:"expiration"`
}

func NewPath(path string, labels map[string]string) string {
	return fmt.Sprintf("%s%s%s", formatPath(path), NAMESPACE_LABEL_SEPARATOR, formatLabels(labels))
}

func parseLabels(key string) map[string]string {
	labels := strings.Split(key, LABEL_SEPARATOR)
	mapLabels := make(map[string]string, len(labels))

	for _, label := range labels {
		key, value := parseLabel(label)
		if len(key) > 0 {
			mapLabels[key] = value
		}
	}
	return mapLabels
}

func parseLabel(label string) (string, string) {
	s := strings.Split(label, LABEL_SEPARATOR)
	return s[0], strings.Join(s[1:], LABEL_SEPARATOR)
}

func formatPath(path string) string {
	if len(path) > 0 {
		path = PATH_SEPARATOR + strings.TrimPrefix(path, PATH_SEPARATOR)
		path = strings.TrimSuffix(path, PATH_SEPARATOR) + PATH_SEPARATOR
	}
	return path
}

func sortLabels(labels map[string]string) []string {
	keyValues := make([]string, len(labels))
	for key, val := range labels {
		keyValues = append(keyValues, fmt.Sprintf("%s%s%s", key, KEY_VALUE_SEPARATOR, val))
	}
	sort.Strings(keyValues)
	return keyValues
}

func formatLabels(labels map[string]string) string {
	keyValues := sortLabels(labels)
	return strings.Join(keyValues, LABEL_SEPARATOR)
}

func FormatPatternKey(path string, labels map[string]string) string {
	return fmt.Sprintf("%s%s", formatPatternPath(path), formatPatternLabels(labels))
}

func formatPatternLabels(labels map[string]string) string {
	keyValues := sortLabels(labels)
	if len(keyValues) == 0 {
		return ""
	}
	return fmt.Sprintf("%s*", strings.Join(keyValues, "*"))
}

func formatPatternPath(path string) string {
	pattern := formatPath(path)
	return pattern + "*"
}
