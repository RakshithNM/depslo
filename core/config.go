package core

// LanguageConfig holds expansion/contraction rates for different languages
type LanguageConfig struct {
	Code          string             `json:"code"`
	Name          string             `json:"name"`
	BaseExpansion float64            `json:"base_expansion"` // Base expansion rate (1.0 = no change)
	LengthFactor  float64            `json:"length_factor"`  // How much length affects expansion
	ContentTypes  map[string]float64 `json:"content_types"`  // Multipliers for different content types
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
				BaseExpansion: 1.25, // 25% expansion
				LengthFactor:  0.8,  // Shorter strings expand more
				ContentTypes: map[string]float64{
					"ui":        1.0,
					"technical": 0.8,
					"marketing": 1.1,
					"legal":     0.9,
				},
			},
			"fr": {
				Code:          "fr",
				Name:          "French",
				BaseExpansion: 1.23,
				LengthFactor:  0.8,
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
				BaseExpansion: 1.30,
				LengthFactor:  0.75,
				ContentTypes: map[string]float64{
					"ui":        1.0,
					"technical": 0.85,
					"marketing": 1.2,
					"legal":     0.9,
				},
			},
			"zh": {
				Code:          "zh",
				Name:          "Chinese",
				BaseExpansion: 0.85, // 15% contraction
				LengthFactor:  1.1,  // Less contraction for shorter strings
				ContentTypes: map[string]float64{
					"ui":        1.0,
					"technical": 1.0,
					"marketing": 0.9,
					"legal":     1.0,
				},
			},
			"ja": {
				Code:          "ja",
				Name:          "Japanese",
				BaseExpansion: 0.90,
				LengthFactor:  1.1,
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
	// Base expansion rate
	rate := lc.BaseExpansion

	// Adjust based on text length (shorter texts expand more)
	lengthAdjustment := 1.0
	if textLength < 10 {
		lengthAdjustment = 1.5 // 50% more expansion for very short strings
	} else if textLength < 50 {
		lengthAdjustment = 1.2 // 20% more expansion for short strings
	} else if textLength > 100 {
		lengthAdjustment = 0.8 // 20% less expansion for long strings
	}

	rate *= lengthAdjustment * lc.LengthFactor

	// Adjust based on content type
	if multiplier, exists := lc.ContentTypes[contentType]; exists {
		rate *= multiplier
	}

	return rate
}
