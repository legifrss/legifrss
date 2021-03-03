package rss

import (
	"time"

	"github.com/gorilla/feeds"
	"github.com/ldicarlo/legifrss/server/pkg/models"
)

type metaRSS struct {
	updated string
}

// TransformToRSS is the Main transformation function
func TransformToRSS(input []models.LegifranceElement, feedDesc models.FeedDescription) *feeds.AtomFeed {
	feed := &feeds.AtomFeed{
		Xmlns:    "http://www.w3.org/2005/Atom",
		Title:    "Legifrance RSS " + feedDesc.TitleSuffix,
		Id:       "https://legifrss.org/",
		Subtitle: "This is a non-official RSS feed for Legifrance's Official Law updates. This is at TESTING stage for now. If you want to follow that topic, you can find more info at https://github.com/ldicarlo/legifrance-rss " + feedDesc.DescriptionSuffix,
		Author:   &feeds.AtomAuthor{AtomPerson: feeds.AtomPerson{Name: "Luca Di Carlo", Email: "luca@di-carlo.fr"}},
		// unnecessary loop; there is another one below
		Updated: extractMeta(input).updated,
		Logo:    "https://www.legifrance.gouv.fr/contenu/logo",
	}

	for _, legifranceElement := range input {
		feed.Entries = append(feed.Entries, transformLegifranceElement(legifranceElement))
	}

	return feed
}

func extractMeta(input []models.LegifranceElement) (meta metaRSS) {
	date := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	for _, element := range input {
		if date.Before(element.Date) {
			date = element.Date
		}
	}
	meta.updated = date.Format(time.RFC3339)
	return
}

func transformLegifranceElement(element models.LegifranceElement) *feeds.AtomEntry {
	return &feeds.AtomEntry{
		Title:     element.Description,
		Links:     []feeds.AtomLink{{Href: "https://www.legifrance.gouv.fr/jorf/id/" + element.ID}},
		Author:    &feeds.AtomAuthor{AtomPerson: feeds.AtomPerson{Name: element.Author}},
		Published: element.Date.Format(time.RFC3339),
		Content:   &feeds.AtomContent{Content: element.Content, Type: "html"},
		Id:        "https://www.legifrance.gouv.fr/jorf/id/" + element.ID,
		Updated:   element.Date.Format(time.RFC3339),
	}
}
