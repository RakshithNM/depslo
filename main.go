package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	core "github.com/RakshithNM/depslo/core"
)

// gin route handler
func translate(c *gin.Context) {
	var inputStrings map[string]string

	if err := c.BindJSON(&inputStrings); err != nil {
		return
	}

	c.IndentedJSON(http.StatusOK, core.PsuedoLocalize(inputStrings))
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
