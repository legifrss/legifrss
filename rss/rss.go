package rss

import (
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/ldicarlo/legifrss/server/models"
)

// TransformToRSS is the Main transformation function
func TransformToRSS(input []models.LegifranceElement, feedDesc models.FeedDescription) *feeds.AtomFeed {
	now := time.Now().String()
	feed := &feeds.AtomFeed{
		Title:    "Legifrance RSS " + feedDesc.TitleSuffix,
		Link:     &feeds.AtomLink{Href: "https://github.com/ldicarlo/legifrance-rss" + feedDesc.LinkSuffix},
		Subtitle: "This is a non-official RSS feed for Legifrance's Official Law updates. This is at TESTING stage for now. If you want to follow that topic, you can find more info at https://github.com/ldicarlo/legifrance-rss" + feedDesc.DescriptionSuffix,
		Author:   &feeds.AtomAuthor{AtomPerson: feeds.AtomPerson{Name: "Luca Di Carlo", Email: "luca@di-carlo.fr"}},
		Updated:  now,
		Category: "French Law",
	}

	for _, legifranceElement := range input {
		feed.Entries = append(feed.Entries, transformLegifranceElement(legifranceElement))
	}

	return feed
}

func transformLegifranceElement(element models.LegifranceElement) *feeds.AtomEntry {
	return &feeds.AtomEntry{
		Title:     "[" + element.Nature + " - " + element.Id + "]: " + element.Description,
		Links:     []feeds.AtomLink{feeds.AtomLink{Href: "https://www.legifrance.gouv.fr/jorf/id/" + element.Id}},
		Category:  strings.Join(element.Category, "/"),
		Author:    &feeds.AtomAuthor{AtomPerson: feeds.AtomPerson{Name: element.Author}},
		Published: element.Date,
		Content:   &feeds.AtomContent{Content: element.Content},
	}
}
