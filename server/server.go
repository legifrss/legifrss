package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ldicarlo/legifrss/server/db"
	"github.com/ldicarlo/legifrss/server/models"
	"github.com/ldicarlo/legifrss/server/rss"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/latest", func(c *gin.Context) {
		queryContext := models.QueryContext{
			Keyword: c.DefaultQuery("q", ""),
			Author:  c.DefaultQuery("author", ""),
			Nature:  c.DefaultQuery("nature", ""),
		}
		result := db.Query(queryContext)
		rss := rss.TransformToRSS(result, models.FeedDescription{})
		c.XML(200, rss)
	})
	r.Run()
}
