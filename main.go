package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	core "github.com/RakshithNM/depslo/core"
)

// LocalizationRequest represents the incoming request structure
type LocalizationRequest struct {
	Strings     map[string]string `json:"strings"`
	Language    string            `json:"language"`     // Target language code
	ContentType string            `json:"content_type"` // ui, technical, marketing, legal
}

// Gin route handler to Pseudo localise passed in JSON string
func translate(c *gin.Context) {
	var inputStrings map[string]string

	if err := c.BindJSON(&inputStrings); err != nil {
		return
	}

	// Pseudo localize the JSON passed in to the endpoint, use the Depslo core pseudo localize function
	c.IndentedJSON(http.StatusOK, core.PseudoLocalize(inputStrings))
}

// Gin route handler to Pseudo localise passed in JSON string
func localize(c *gin.Context) {
	var req LocalizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pseudoStrings := core.PseudoLocalizeAdvanced(req.Strings, req.Language, req.ContentType)
	c.JSON(http.StatusOK, pseudoStrings)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/translate", translate)
	router.POST("/localize", localize)

	router.Run("localhost:1234")
}
