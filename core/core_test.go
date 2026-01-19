package core

import (
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
	target := int(float64(srcLen) * rate)
	gotLen := utf8.RuneCountInString(got)

	if gotLen != target {
		t.Fatalf("expected length %d, got %d", target, gotLen)
	}
}
