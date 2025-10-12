package service

import (
	"errors"
	"strings"

	"github.com/Sevacoming/sprint6/pkg/morse"
)

var ErrEmptyInput = errors.New("empty input")

// DetectAndConvert определяет формат входных данных и конвертирует:
// - Морзе -> текст
// - Текст -> Морзе
func DetectAndConvert(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", ErrEmptyInput
	}
	if isLikelyMorse(trimmed) {
		return morse.ToText(trimmed), nil
	}
	return morse.ToMorse(trimmed), nil
}

// допускаем только . - / и пробельные — тогда это похоже на Морзе
func isLikelyMorse(s string) bool {
	for _, r := range s {
		switch r {
		case '.', '-', '/', ' ', '\n', '\r', '\t':
			// ok
		default:
			return false
		}
	}
	return true
}
