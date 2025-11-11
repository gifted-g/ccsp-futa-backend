package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListEventsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"events": []string{}})
}
func ListContactsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"contacts": []string{}})
}
