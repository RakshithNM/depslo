package core

import "math"

type LanguageConfig struct {
	Code          string             `json:"code"`
	Name          string             `json:"name"`
	BaseExpansion float64            `json:"base_expansion"` // neutral medium length factor (1.0 = same length, 1.25 = +25%)
	ShortBonus    float64            `json:"short_bonus"`    // multiplier strength applied for short strings (as a % of base)
	LongPenalty   float64            `json:"long_penalty"`   // multiplier strength applied for long strings (as a % of base)
	MinRate       float64            `json:"min_rate"`       // language-specific floor
	MaxRate       float64            `json:"max_rate"`       // language-specific ceiling
	ContentTypes  map[string]float64 `json:"content_types"`  // extra multiplier per content-type
}

// PseudoLocalizationConfig holds all language configurations
type PseudoLocalizationConfig struct {
	Languages map[string]LanguageConfig `json:"languages"`
}

// GetDefaultConfig tuned to keep languages distinct
func GetDefaultConfig() PseudoLocalizationConfig {
	return PseudoLocalizationConfig{
		Languages: map[string]LanguageConfig{
			"es": {
				Code:          "es",
				Name:          "Spanish",
				BaseExpansion: 1.30, // medium strings trend around +30% vs English
				ShortBonus:    0.55, // short UI strings can approach/exceed 2x
				LongPenalty:   0.12, // long strings expand less
				MinRate:       1.05,
				MaxRate:       2.10,
				ContentTypes: map[string]float64{
					"ui":        1.00,
					"technical": 1.05,
					"marketing": 1.10,
					"legal":     1.15,
				},
			},
			"de": {
				Code:          "de",
				Name:          "German",
				BaseExpansion: 1.35, // medium strings trend around +30% to +35%
				ShortBonus:    0.55,
				LongPenalty:   0.12,
				MinRate:       1.10,
				MaxRate:       2.20,
				ContentTypes: map[string]float64{
					"ui":        1.00,
					"technical": 1.08,
					"marketing": 1.12,
					"legal":     1.18,
				},
			},
			"fr": {
				Code:          "fr",
				Name:          "French",
				BaseExpansion: 1.30, // medium strings trend around +30% vs English
				ShortBonus:    0.55,
				LongPenalty:   0.12,
				MinRate:       1.05,
				MaxRate:       2.10,
				ContentTypes: map[string]float64{
					"ui":        1.00,
					"technical": 1.05,
					"marketing": 1.10,
					"legal":     1.15,
				},
			},
			"zh": {
				Code:          "zh",
				Name:          "Chinese",
				BaseExpansion: 0.88, // often contracts vs English, but not universally
				ShortBonus:    0.25,
				LongPenalty:   0.10,
				MinRate:       0.65,
				MaxRate:       1.15,
				ContentTypes: map[string]float64{
					"ui":        1.00,
					"technical": 0.98,
					"marketing": 0.92,
					"legal":     1.03,
				},
			},
		},
	}
}

// lengthWeight returns a value in [-1..1] describing how “short” or “long” the string is.
//
//	+1.0  = very short (<=5 chars)
//	 0.0  = medium (31..100 chars)
//	-1.0  = very long (> 100 chars)
func lengthWeight(n int) float64 {
	switch {
	case n <= 5:
		return 1.0
	case n <= 15:
		return 0.7
	case n <= 30:
		return 0.3
	case n <= 100:
		return 0.0
	default:
		return -1.0
	}
}

// CalculateExpansionRate calculates a language-distinct expansion rate.
func (lc LanguageConfig) CalculateExpansionRate(textLength int, contentType string) float64 {
	// Base multiplier for medium strings
	rate := lc.BaseExpansion

	// Length-driven relative adjustment:
	// Positive weight (short strings) amplifies ShortBonus;
	// Negative weight (long strings) amplifies LongPenalty.
	w := lengthWeight(textLength)
	if w > 0 {
		rate *= (1.0 + lc.ShortBonus*w)
	} else if w < 0 {
		rate *= (1.0 + lc.LongPenalty*w) // w is negative => reduces rate
	}

	// Content-type multiplier (kept language-specific)
	if mult, ok := lc.ContentTypes[contentType]; ok {
		rate *= mult
	}

	// Clamp per-language to keep behavior believable
	if lc.MinRate > 0 {
		rate = math.Max(lc.MinRate, rate)
	}
	if lc.MaxRate > 0 {
		rate = math.Min(lc.MaxRate, rate)
	}
	return rate
}
