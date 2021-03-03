package main

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ldicarlo/legifrss/server/pkg/db"
	"github.com/ldicarlo/legifrss/server/pkg/models"
	"github.com/ldicarlo/legifrss/server/pkg/rss"

	cache "github.com/stockx/go-gin-cache"
	"github.com/stockx/go-gin-cache/persistence"
)

func main() {

	store := persistence.NewInMemoryStore(time.Second)
	r := gin.Default()
	r.GET("/latest", func(c *gin.Context) {
		queryContext := models.QueryContext{
			Keyword: strings.ToUpper(c.DefaultQuery("q", "")),
			Author:  strings.ToUpper(c.DefaultQuery("author", "")),
			Nature:  strings.ToUpper(c.DefaultQuery("nature", "")),
		}

		result := db.Query(queryContext)
		rss := rss.TransformToRSS(result, models.FeedDescription{
			TitleSuffix:       strings.TrimSpace(queryContext.Author + " " + queryContext.Nature + " " + queryContext.Keyword),
			DescriptionSuffix: "",
			LinkSuffix:        c.Request.Host + c.FullPath() + c.Request.RequestURI,
		})
		c.XML(200, rss)
	})
	r.GET("/authors", cache.CachePage(store, time.Minute, func(c *gin.Context) { c.JSON(200, db.GetAuthors()) }))
	r.GET("/natures", cache.CachePage(store, time.Minute, func(c *gin.Context) { c.JSON(200, db.GetNatures()) }))
	r.Run(":8080")
}
