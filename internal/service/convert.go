package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

var (
	runeToMorse = map[rune]string{
		'А': ".-", 'Б': "-...", 'В': ".--", 'Г': "--.", 'Д': "-..",
		'Е': ".", 'Ж': "...-", 'З': "--..", 'И': "..", 'Й': ".---",
		'К': "-.-", 'Л': ".-..", 'М': "--", 'Н': "-.", 'О': "---",
		'П': ".--.", 'Р': ".-.", 'С': "...", 'Т': "-", 'У': "..-",
		'Ф': "..-.", 'Х': "....", 'Ц': "-.-.", 'Ч': "---.", 'Ш': "----",
		'Щ': "--.-", 'Ы': "-.--", 'Ь': "-..-", 'Э': "..-..", 'Ю': "..--",
		'Я': ".-.-",
		'0': "-----", '1': ".----", '2': "..---", '3': "...--", '4': "....-",
		'5': ".....", '6': "-....", '7': "--...", '8': "---..", '9': "----.",
	}
	morseToRune map[string]rune
)

func init() {
	morseToRune = make(map[string]rune, len(runeToMorse))
	for r, m := range runeToMorse {
		morseToRune[m] = r
	}
}

// Convert — принимает СТРОКУ, автоопределяет Морзе или обычный текст и конвертирует.
func Convert(input string) (string, error) {
	if isMorse(input) {
		return morseToText(input)
	}
	return textToMorse(input)
}

// ConvertFromReader — адаптер, если у тебя io.Reader (multipart.File и т.п.).
func ConvertFromReader(r io.Reader) (string, error) {
	if r == nil {
		return "", errors.New("nil reader")
	}
	var sb strings.Builder
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if sb.Len() > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(sc.Text())
	}
	if err := sc.Err(); err != nil {
		return "", err
	}
	return Convert(sb.String())
}

// ------ helpers ------

func isMorse(s string) bool {
	for _, r := range s {
		switch r {
		case '.', '-', '/':
			continue
		default:
			if unicode.IsSpace(r) {
				continue
			}
			return false
		}
	}
	return strings.TrimSpace(s) != ""
}

func textToMorse(s string) (string, error) {
	var b strings.Builder
	needSpace := false

	writeLetter := func(r rune) error {
		if r == 'Ё' {
			r = 'Е'
		}
		code, ok := runeToMorse[r]
		if !ok {
			return fmt.Errorf("неизвестный символ для Морзе: %q", r)
		}
		if needSpace {
			b.WriteByte(' ')
		}
		b.WriteString(code)
		needSpace = true
		return nil
	}

	for _, r := range strings.ToUpper(s) {
		switch r {
		case '\n':
			b.WriteByte('\n')
			needSpace = false
		case '\r', '\t':
		case ' ':
			if needSpace {
				b.WriteString(" /")
			} else {
				b.WriteString("/")
			}
			needSpace = true
		default:
			if err := writeLetter(r); err != nil {
				return "", err
			}
		}
	}
	out := strings.TrimSpace(b.String())
	if out == "" {
		return "", errors.New("пустой вход")
	}
	return out, nil
}

func morseToText(s string) (string, error) {
	lineNorm := func(line string) string {
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "   ", " / ")
		return line
	}

	var out strings.Builder
	lines := strings.Split(s, "\n")
	for li, line := range lines {
		line = lineNorm(line)
		if line == "" {
			if li > 0 {
				out.WriteByte('\n')
			}
			continue
		}
		tokens := strings.Fields(line)
		for _, t := range tokens {
			if t == "/" {
				out.WriteByte(' ')
				continue
			}
			r, ok := morseToRune[t]
			if !ok {
				return "", fmt.Errorf("неизвестная последовательность Морзе: %q", t)
			}
			out.WriteRune(r)
		}
		if li < len(lines)-1 {
			out.WriteByte('\n')
		}
	}
	result := out.String()
	if strings.TrimSpace(result) == "" {
		return "", errors.New("пустой результат декодирования")
	}
	return result, nil
}
