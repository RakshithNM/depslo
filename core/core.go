/*
 * DEPSLO core
 */
package core

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Propose a length for the psuedoLocalization of string
func proposeLength(inString string) int {
	// To read values from LENGTHINCREASEMAP in order
	keys := make([]int, 0)
	for k, _ := range LENGTHINCREASEMAP {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	length := len(inString)
	for _, key := range keys {
		if length > key {
			continue
		}
		return LENGTHINCREASEMAP[key][1]
	}
	// proposing longest length
	// TODO: maybe should be randomised to return a number within the range(pending research)
	return LENGTHINCREASEMAP[10][1]
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

// PsuedoLocalize the JSON
func PsuedoLocalize(inJSON map[string]string) map[string]string {
	for i := range inJSON {
		stringToTranslate := inJSON[i]
		proposedLength := proposeLength(stringToTranslate)
		var translatedBuffer bytes.Buffer
		for _, s := range stringToTranslate {
			key := rune(s)
			_, isPresent := LETTERS[key]
			if isPresent {
				translatedBuffer.WriteRune(LETTERS[key])
			}
		}
		var elongatedString string
		if len(translatedBuffer.String()) > 0 {
			translatedString := translatedBuffer.String()
			elongatedString = elongateToLength(translatedString, proposedLength)
		} else {
			elongatedString = stringToTranslate
		}
		inJSON[i] = elongatedString
	}
	return inJSON
}
