package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	core "github.com/RakshithNM/depslo/core"
)

// Gin route handler to PsuedoLocalise passsed in JSON string
func translate(c *gin.Context) {
	var inputStrings map[string]string

	if err := c.BindJSON(&inputStrings); err != nil {
		return
	}

	// PsuedoLocalize the JSON passed in to the endpoint, use the depslo core PsuedoLocalize function
	c.IndentedJSON(http.StatusOK, core.PsuedoLocalize(inputStrings))
}

// Ping to check if server is running
func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "pong")
}

func main() {
	router := gin.Default()
	router.GET("/ping", ping)
	router.POST("/translate", translate)

	router.Run("localhost:1234")
}
