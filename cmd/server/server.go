package main

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"legifrss/pkg/db"
	"legifrss/pkg/models"
	"legifrss/pkg/rss"

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
			LinkSuffix:        c.Request.Host + c.Request.RequestURI,
		})
		c.XML(200, rss)
	})
	r.GET("/authors", cache.CachePage(store, time.Minute, func(c *gin.Context) { c.JSON(200, db.GetAuthors()) }))
	r.GET("/natures", cache.CachePage(store, time.Minute, func(c *gin.Context) { c.JSON(200, db.GetNatures()) }))

	// r.GET("/callback", func(c *gin.Context) {
	// 	bot.RegisterToken(c.Query("oauth_token"), c.Query("oauth_verifier"))
	// 	c.JSON(200, "OK")
	// })
	// r.GET("/auth", func(c *gin.Context) { c.Redirect(307, bot.GetAuthURL()) })
	// For local purposes only
	// r.GET("/twitter_publish", func(c *gin.Context) {
	// 	bot.ProcessElems()
	// 	c.JSON(200, "ok")
	// })
	// r.GET("/twitter_statuses", func(c *gin.Context) {
	// 	c.JSON(200, bot.GetAllTweets())
	// })
	// r.GET("/remove_twitter_statuses", func(c *gin.Context) {
	// 	c.JSON(200, bot.RemoveAllTweets())
	// })

	r.Run(":8080")
}
