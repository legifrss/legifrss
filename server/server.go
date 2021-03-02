package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ldicarlo/legifrss/server/db"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/latest", func(c *gin.Context) {
		keyword := c.DefaultQuery("q", "")
		println(keyword)
		c.XML(200, db.GetAll())
	})
	r.Run()
}
