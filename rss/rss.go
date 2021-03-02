package rss

import (
	"time"

	"github.com/gorilla/feeds"
	"github.com/ldicarlo/legifrss/server/models"
)

// TransformToRSS is the Main transformation function
func TransformToRSS(input []models.LegifranceElement, feedDesc models.FeedDescription) *feeds.AtomFeed {
	now := time.Now().Format(time.RFC3339)
	feed := &feeds.AtomFeed{
		Xmlns:    "http://www.w3.org/2005/Atom",
		Title:    "Legifrance RSS " + feedDesc.TitleSuffix,
		Id:       "https://raw.githubusercontent.com/ldicarlo/legifrance-rss/master/feed/" + feedDesc.LinkSuffix,
		Subtitle: "This is a non-official RSS feed for Legifrance's Official Law updates. This is at TESTING stage for now. If you want to follow that topic, you can find more info at https://github.com/ldicarlo/legifrance-rss " + feedDesc.DescriptionSuffix,
		Author:   &feeds.AtomAuthor{AtomPerson: feeds.AtomPerson{Name: "Luca Di Carlo", Email: "luca@di-carlo.fr"}},
		Updated:  now,
		Logo:     "https://www.legifrance.gouv.fr/contenu/logo",
	}

	for _, legifranceElement := range input {
		feed.Entries = append(feed.Entries, transformLegifranceElement(legifranceElement, now))
	}

	return feed
}

func AggregateRSS(feedDesc *feeds.AtomFeed, feedEntries []*feeds.AtomEntry) {
	return true
}

func transformLegifranceElement(element models.LegifranceElement, date string) *feeds.AtomEntry {
	return &feeds.AtomEntry{
		Title:     element.Description,
		Links:     []feeds.AtomLink{{Href: "https://www.legifrance.gouv.fr/jorf/id/" + element.Id}},
		Author:    &feeds.AtomAuthor{AtomPerson: feeds.AtomPerson{Name: element.Author}},
		Published: element.Date.Format(time.RFC3339),
		Content:   &feeds.AtomContent{Content: element.Content, Type: "html"},
		Id:        "https://www.legifrance.gouv.fr/jorf/id/" + element.Id,
		Updated:   date,
	}
}
