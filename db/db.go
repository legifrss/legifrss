package db

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gorilla/feeds"
	"github.com/ldicarlo/legifrss/server/models"
	"github.com/ldicarlo/legifrss/server/rss"
	"github.com/ldicarlo/legifrss/server/utils"
)

// GetAll returns all contents of db.xml
func GetAll() (feed feeds.AtomFeed) {
	file, err := os.Open("feed/db.xml")
	utils.ErrCheck(err)

	byteValue, _ := ioutil.ReadAll(file)

	xml.Unmarshal(byteValue, &feed)
	return feed
}

func Query(queryContext models.QueryContext) feeds.AtomFeed {
	var feed = GetAll()
	var entries []*feeds.AtomEntry

	for _, element := range feed.Entries {
		if keep(queryContext, element) {
			entries = append(entries, element)
		}
	}

	return rss.TransformToRSS()
}

func keep(queryContext models.QueryContext, element *feeds.AtomEntry) bool {
	return strings.Contains(element.Content.Content, queryContext.Keyword)

}
