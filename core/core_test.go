package core

import (
	"math"
	"testing"
	"unicode/utf8"
)

func TestPseudoLocalizeDefaultExpansion(t *testing.T) {
	input := map[string]string{
		"TITLE": "Hello world",
	}
	output := PseudoLocalize(input)
	got := output["TITLE"]
	if got == "" {
		t.Fatalf("expected non-empty output")
	}

	cfg := GetDefaultConfig()
	langCfg := cfg.Languages["es"]
	srcLen := utf8.RuneCountInString(input["TITLE"])
	rate := langCfg.CalculateExpansionRate(srcLen, "ui")
	rawTarget := float64(srcLen) * rate
	target := int(math.Floor(rawTarget))
	if rate >= 1.0 {
		target = int(math.Ceil(rawTarget))
	}
	if srcLen > 0 && target == 0 {
		target = 1
	}
	gotLen := utf8.RuneCountInString(got)

	if gotLen != target {
		t.Fatalf("expected length %d, got %d", target, gotLen)
	}
}

func TestGeneratePseudoTextRoundsUpDuringExpansion(t *testing.T) {
	got := generatePseudoText("AB", 1.30)
	if utf8.RuneCountInString(got) != 3 {
		t.Fatalf("expected ceil behavior for expansion: length 3, got %d", utf8.RuneCountInString(got))
	}
}

func TestGeneratePseudoTextKeepsNonEmptyForNonEmptyInput(t *testing.T) {
	got := generatePseudoText("A", 0.60)
	if utf8.RuneCountInString(got) != 1 {
		t.Fatalf("expected non-empty output for non-empty input, got length %d", utf8.RuneCountInString(got))
	}
}
