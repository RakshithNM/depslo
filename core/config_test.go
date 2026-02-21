package core

import "testing"

func TestExpansionRateRanges(t *testing.T) {
	cfg := GetDefaultConfig()

	tests := []struct {
		lang       string
		shortLen   int
		mediumLen  int
		longLen    int
		expectMin  float64
		expectMax  float64
		contracted bool
	}{
		{"es", 5, 40, 120, 1.05, 2.10, false},
		{"de", 5, 40, 120, 1.10, 2.20, false},
		{"fr", 5, 40, 120, 1.05, 2.10, false},
		{"zh", 5, 40, 120, 0.65, 1.15, true},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			langCfg, ok := cfg.Languages[tt.lang]
			if !ok {
				t.Fatalf("missing language config for %s", tt.lang)
			}

			shortRate := langCfg.CalculateExpansionRate(tt.shortLen, "ui")
			mediumRate := langCfg.CalculateExpansionRate(tt.mediumLen, "ui")
			longRate := langCfg.CalculateExpansionRate(tt.longLen, "ui")

			for _, r := range []float64{shortRate, mediumRate, longRate} {
				if r < langCfg.MinRate || r > langCfg.MaxRate {
					t.Fatalf("rate out of bounds for %s: %f", tt.lang, r)
				}
				if r < tt.expectMin || r > tt.expectMax {
					t.Fatalf("rate outside expected range for %s: %f", tt.lang, r)
				}
			}

			if !(shortRate >= mediumRate && mediumRate >= longRate) {
				t.Fatalf("expected short >= medium >= long for %s; got %f %f %f", tt.lang, shortRate, mediumRate, longRate)
			}

			if tt.contracted {
				if mediumRate > 1.0 {
					t.Fatalf("expected contraction for %s; got %f", tt.lang, mediumRate)
				}
			} else {
				if mediumRate < 1.0 {
					t.Fatalf("expected expansion for %s; got %f", tt.lang, mediumRate)
				}
			}
		})
	}
}

func TestContentTypeMultipliers(t *testing.T) {
	cfg := GetDefaultConfig()

	expandLangs := []string{"es", "de", "fr"}
	for _, lang := range expandLangs {
		langCfg := cfg.Languages[lang]
		uiRate := langCfg.CalculateExpansionRate(40, "ui")
		marketingRate := langCfg.CalculateExpansionRate(40, "marketing")
		if marketingRate <= uiRate {
			t.Fatalf("expected marketing to expand more than ui for %s", lang)
		}
	}

	zhCfg := cfg.Languages["zh"]
	zhUiRate := zhCfg.CalculateExpansionRate(40, "ui")
	zhMarketingRate := zhCfg.CalculateExpansionRate(40, "marketing")
	if zhMarketingRate >= zhUiRate {
		t.Fatalf("expected marketing to contract more than ui for zh")
	}
}

func TestPracticalRuleAlignment(t *testing.T) {
	cfg := GetDefaultConfig()

	european := []string{"es", "de", "fr"}
	mediumLen := 40
	shortLen := 5

	avgMediumRate := 0.0
	upToFiftyCount := 0

	for _, lang := range european {
		langCfg := cfg.Languages[lang]
		mediumRate := langCfg.CalculateExpansionRate(mediumLen, "ui")
		shortRate := langCfg.CalculateExpansionRate(shortLen, "ui")
		legalRate := langCfg.CalculateExpansionRate(mediumLen, "legal")

		avgMediumRate += mediumRate

		if shortRate < 1.95 {
			t.Fatalf("expected short-string stress near 2x for %s, got %f", lang, shortRate)
		}

		if langCfg.MaxRate < 2.0 {
			t.Fatalf("expected max rate >= 2.0 for %s, got %f", lang, langCfg.MaxRate)
		}

		if legalRate >= 1.49 {
			upToFiftyCount++
		}
	}

	avgMediumRate /= float64(len(european))
	if avgMediumRate < 1.28 || avgMediumRate > 1.40 {
		t.Fatalf("expected european medium-rate average near +30%%, got %f", avgMediumRate)
	}

	if upToFiftyCount < 2 {
		t.Fatalf("expected many european languages to approach +50%% in legal content, got %d", upToFiftyCount)
	}

	zhCfg := cfg.Languages["zh"]
	zhMediumUI := zhCfg.CalculateExpansionRate(mediumLen, "ui")
	if zhMediumUI >= 1.0 {
		t.Fatalf("expected zh medium ui to typically contract, got %f", zhMediumUI)
	}

	zhShortLegal := zhCfg.CalculateExpansionRate(shortLen, "legal")
	if zhShortLegal <= 1.0 {
		t.Fatalf("expected zh to allow occasional expansion for short/legal text, got %f", zhShortLegal)
	}
}
