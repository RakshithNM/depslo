package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
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
	return LENGTHINCREASEMAP[10][1]
}

// Elongate the string to the desired length
func elongateToLength(inString string, inLength int) string {
	expectedLength := inLength
	currentLength := len(inString)
	var localElongatedString string
	count := 1
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

// psuedoLocalize the JSON
func psuedoLocalize(inJSON map[string]interface{}) map[string]interface{} {
	for i := range inJSON {
		stringToTranslate := inJSON[i].(string)
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

// gin route handler
func translate(c *gin.Context) {
	var inputStrings map[string]interface{}

	if err := c.BindJSON(&inputStrings); err != nil {
		return
	}

	c.IndentedJSON(http.StatusOK, psuedoLocalize(inputStrings))
}

// ping to check if server is running
func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "pong")
}

func main() {
	router := gin.Default()
	router.GET("/ping", ping)
	router.POST("/translate", translate)

	router.Run("localhost:1234")
}
