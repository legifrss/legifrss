package rss

import (
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/ldicarlo/legifrss/server/models"
)

// TransformToRSS is the Main transformation function
func TransformToRSS(input []models.LegifranceElement, feedDesc models.FeedDescription) *feeds.RssFeed {
	now := time.Now().String()
	feed := &feeds.RssFeed{
		Title:          "Legifrance RSS" + feedDesc.TitleSuffix,
		Link:           "https://github.com/ldicarlo/legifrance-rss" + feedDesc.LinkSuffix,
		Description:    "This is a non-official RSS feed for Legifrance's Official Law updates. This is at TESTING stage for now. If you want to follow that topic, you can find more info at https://github.com/ldicarlo/legifrss" + feedDesc.DescriptionSuffix,
		ManagingEditor: "luca@di-carlo.fr (Luca Di Carlo)",
		LastBuildDate:  now,
		PubDate:        now,
	}

	for _, legifranceElement := range input {
		feed.Items = append(feed.Items, transformLegifranceElement(legifranceElement))
	}

	return feed
}

func transformLegifranceElement(element models.LegifranceElement) *feeds.RssItem {
	return &feeds.RssItem{
		Title:       element.Nature + " - " + element.Id,
		Link:        "https://www.legifrance.gouv.fr/jorf/id/" + element.Id,
		Description: element.Title,
		Category:    strings.Join(element.Category, "/"),
		Author:      element.Author,
	}
}
