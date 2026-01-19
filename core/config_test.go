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
		{"es", 5, 40, 120, 1.05, 1.75, false},
		{"de", 5, 40, 120, 1.10, 1.85, false},
		{"fr", 5, 40, 120, 1.05, 1.70, false},
		{"zh", 5, 40, 120, 0.60, 1.00, true},
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
