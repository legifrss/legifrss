package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ldicarlo/legifrss/server/pkg/db"
	"github.com/ldicarlo/legifrss/server/pkg/models"
	"github.com/ldicarlo/legifrss/server/pkg/rss"
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
			Keyword: strings.ToUpper(c.DefaultQuery("q", "")),
			Author:  strings.ToUpper(c.DefaultQuery("author", "")),
			Nature:  strings.ToUpper(c.DefaultQuery("nature", "")),
		}

		result := db.Query(queryContext)
		rss := rss.TransformToRSS(result, models.FeedDescription{})
		c.XML(200, rss)
	})
	// TODO add handler for natures and authors
	r.Run()
}
