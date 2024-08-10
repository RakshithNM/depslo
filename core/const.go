/* DEPSLO core
 *
 * All the constants used in depslo source code are in this file
 *
 * HELP NEEDED: LENGTHINCREASEMAP should be bettered altered to map the reality of string elongation in other languages in comparision to the english language.
 * If there are any resources on research done on this subject, please let me know. You are also welcome to submit a PR.
 */
package core

import "math"

// VOWELS slice has all the vowels in English language
var VOWELS = []rune{'a', 'e', 'i', 'o', 'u', 'y', 'A', 'E', 'I', 'O', 'U', 'Y'}

// LETTERS map has all the psuedo localised characters for characters in English language
var LETTERS = map[rune]rune{
	'a': 'α',
	'b': 'ḅ',
	'c': 'ͼ',
	'd': 'ḍ',
	'e': 'ḛ',
	'f': 'ϝ',
	'g': 'ḡ',
	'h': 'ḥ',
	'i': 'ḭ',
	'j': 'ĵ',
	'k': 'ḳ',
	'l': 'ḽ',
	'm': 'ṃ',
	'n': 'ṇ',
	'o': 'ṓ',
	'p': 'ṗ',
	'q': 'ʠ',
	'r': 'ṛ',
	's': 'ṡ',
	't': 'ṭ',
	'u': 'ṵ',
	'v': 'ṽ',
	'w': 'ẁ',
	'x': 'ẋ',
	'y': 'ẏ',
	'z': 'ẓ',
	'A': 'Ḁ',
	'B': 'Ḃ',
	'C': 'Ḉ',
	'D': 'Ḍ',
	'E': 'Ḛ',
	'F': 'Ḟ',
	'G': 'Ḡ',
	'H': 'Ḥ',
	'I': 'Ḭ',
	'J': 'Ĵ',
	'K': 'Ḱ',
	'L': 'Ḻ',
	'M': 'Ṁ',
	'N': 'Ṅ',
	'O': 'Ṏ',
	'P': 'Ṕ',
	'Q': 'Ǫ',
	'R': 'Ṛ',
	'S': 'Ṣ',
	'T': 'Ṫ',
	'U': 'Ṳ',
	'V': 'Ṿ',
	'W': 'Ŵ',
	'X': 'Ẋ',
	'Y': 'Ŷ',
	'Z': 'Ż',
}

// LENGTHINCREASEMAP map has suggested length elongation for psudeo localised strings in other languages compared to English lanugage
var LENGTHINCREASEMAP = map[int][2]int{
	10:             [2]int{200, 300},
	20:             [2]int{180, 200},
	30:             [2]int{160, 180},
	50:             [2]int{140, 160},
	70:             [2]int{151, 170},
	math.MaxUint32: [2]int{130, 130},
}
