package core

import "math"

// LanguageConfig holds expansion/contraction rates for different languages
type LanguageConfig struct {
	Code          string             `json:"code"`
	Name          string             `json:"name"`
	BaseExpansion float64            `json:"base_expansion"`
	ShortBonus    float64            `json:"short_bonus"`
	LongPenalty   float64            `json:"long_penalty"`
	ContentTypes  map[string]float64 `json:"content_types"`
}

// PseudoLocalizationConfig holds all language configurations
type PseudoLocalizationConfig struct {
	Languages map[string]LanguageConfig `json:"languages"`
}

// GetDefaultConfig is the configuration based on research
func GetDefaultConfig() PseudoLocalizationConfig {
	return PseudoLocalizationConfig{
		Languages: map[string]LanguageConfig{
			"es": {
				Code:          "es",
				Name:          "Spanish",
				BaseExpansion: 1.25, // 25% expansion for medium strings
				ShortBonus:    0.75, // Up to 75% additional for very short strings
				LongPenalty:   0.10, // 10% less for long strings
				ContentTypes: map[string]float64{
					"ui":        1.0,
					"technical": 0.8,
					"marketing": 1.1,
					"legal":     0.9,
				},
			},
			"de": {
				Code:          "de",
				Name:          "German",
				BaseExpansion: 1.30, // 30% expansion for medium strings
				ShortBonus:    1.00, // Up to 100% additional for very short strings
				LongPenalty:   0.15, // 15% less for long strings
				ContentTypes: map[string]float64{
					"ui":        1.0,
					"technical": 0.85,
					"marketing": 1.2,
					"legal":     0.9,
				},
			},
			"fr": {
				Code:          "fr",
				Name:          "French",
				BaseExpansion: 1.23, // 23% expansion for medium strings
				ShortBonus:    0.80, // Up to 80% additional for very short strings
				LongPenalty:   0.08, // 8% less for long strings
				ContentTypes: map[string]float64{
					"ui":        1.0,
					"technical": 0.8,
					"marketing": 1.1,
					"legal":     0.9,
				},
			},
			"zh": {
				Code:          "zh",
				Name:          "Chinese",
				BaseExpansion: 0.85, // 15% contraction for medium strings
				ShortBonus:    0.25, // Up to 25% additional (less contraction) for short strings
				LongPenalty:   0.05, // 5% more contraction for long strings
				ContentTypes: map[string]float64{
					"ui":        1.0,
					"technical": 1.0,
					"marketing": 0.9,
					"legal":     1.0,
				},
			},
		},
	}
}

// CalculateExpansionRate calculates the expansion rate for a given string
func (lc LanguageConfig) CalculateExpansionRate(textLength int, contentType string) float64 {
	// Start with base expansion rate
	rate := lc.BaseExpansion

	// Apply length-based adjustments
	if textLength <= 5 {
		// Very short strings (1-5 chars) get maximum bonus
		rate += lc.ShortBonus
	} else if textLength <= 15 {
		// Short strings (6-15 chars) get partial bonus
		rate += lc.ShortBonus * 0.7
	} else if textLength <= 30 {
		// Medium-short strings (16-30 chars) get small bonus
		rate += lc.ShortBonus * 0.3
	} else if textLength > 100 {
		// Long strings get penalty
		rate -= lc.LongPenalty
	}
	// Medium strings (31-100 chars) use base rate

	// Apply content type multiplier
	if multiplier, exists := lc.ContentTypes[contentType]; exists {
		rate *= multiplier
	}

	// Ensure minimum expansion rate
	return math.Max(0.8, rate)
}
