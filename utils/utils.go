package utils

import "github.com/gin-gonic/gin"

func ErrorResponse(c *gin.Context, err error, statuscode int) {
	c.JSON(statuscode, gin.H{"error:": err.Error()})
	return
}
