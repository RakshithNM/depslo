// Package core
package core

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

// transform to exact target length by trimming or padding (rune-aware)
func transformTextToTargetLength(s string, targetLen int) string {
	r := []rune(s)
	if len(r) == targetLen {
		return s
	}
	if len(r) > targetLen {
		return string(r[:targetLen]) // contract
	}
	// expand
	paddingChar := '⬤'
	pad := make([]rune, targetLen-len(r))
	for i := range pad {
		pad[i] = paddingChar
	}
	return s + string(pad)
}

// Advanced pseudo localized text generation
func generatePseudoText(originalText string, expansionRate float64) string {
	if originalText == "" {
		return ""
	}
	// map characters
	var b strings.Builder
	for _, ch := range originalText {
		if val, ok := LETTERS[ch]; ok {
			b.WriteRune(val)
		} else {
			b.WriteRune(ch)
		}
	}
	translated := b.String()

	// compute target (rune-based) and adjust both ways
	target := int(float64(utf8.RuneCountInString(originalText)) * expansionRate)
	if target < 0 {
		target = 0
	}
	return transformTextToTargetLength(translated, target)
}

// PseudoLocalizeAdvanced - Advanced pseudo localization function
func PseudoLocalizeAdvanced(inJSON map[string]string, inLanguage, inContentType string) map[string]string {
	// defaults
	language := "es"
	if inLanguage != "" {
		language = strings.ToLower(inLanguage)
	}
	contentType := "ui"
	if inContentType != "" {
		contentType = strings.ToLower(inContentType)
	}

	cfg := GetDefaultConfig()

	// try chosen language, fallback to ES if missing
	langConfig, ok := cfg.Languages[language]
	if !ok {
		if fallback, ok2 := cfg.Languages["es"]; ok2 {
			langConfig = fallback
		} else {
			// last resort: just return input unchanged
			out := make(map[string]string, len(inJSON))
			for k, v := range inJSON { out[k] = v }
			return out
		}
	}

	out := make(map[string]string, len(inJSON))
	for key, text := range inJSON {
		rate := langConfig.CalculateExpansionRate(utf8.RuneCountInString(text), contentType)
		out[key] = generatePseudoText(text, rate) // see #2
	}
	return out
}

// Propose a length for the pseudo localization of string
func proposeLength(s string) int {
	// LENGTHINCREASEMAP keys presumably refer to original length “buckets”.
	// Use rune count, not bytes.
	n := utf8.RuneCountInString(s)

	keys := make([]int, 0, len(LENGTHINCREASEMAP))
	for k := range LENGTHINCREASEMAP { keys = append(keys, k) }
	sort.Ints(keys)

	for _, key := range keys {
		if n <= key {
			return LENGTHINCREASEMAP[key][1]
		}
	}
	return LENGTHINCREASEMAP[10][1] // your existing “longest” bucket
}

// Elongate the string to the desired length
func elongateToLength(inString string, inLength int) string {
	expectedLength := inLength
	currentLength := len(inString)
	var localElongatedString string
	count := 1
	// 💩
	if currentLength == 0 {
		fmt.Println("ERROR: Empty string, nothing to do!")
		return inString
	}
	for currentLength < expectedLength {
		count += 1
		localElongatedString = strings.Repeat(inString, count)
		currentLength = len(localElongatedString)
	}
	return localElongatedString
}

// PseudoLocalize the JSON
func PseudoLocalize(inJSON map[string]string) map[string]string {
	cfg := GetDefaultConfig()
	langConfig, ok := cfg.Languages["es"]
	if !ok {
		langConfig = LanguageConfig{BaseExpansion: 1.25, MinRate: 1.0, MaxRate: 2.0}
	}

	out := make(map[string]string, len(inJSON))
	for k, v := range inJSON {
		rate := langConfig.CalculateExpansionRate(utf8.RuneCountInString(v), "ui")
		out[k] = generatePseudoText(v, rate)
	}
	return out
}
