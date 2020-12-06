package rss

import (
	"log"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/ldicarlo/legifrss/server/dila"
)

// TransformToRSS is the Main transformation function
func TransformToRSS(result []dila.JOContainerResult) string {
	now := time.Now().String()
	feed := &feeds.RssFeed{
		Title:          "Legifrance RSS",
		Link:           "https://github.com/ldicarlo/legifrance-rss",
		Description:    "This is a non-official RSS feed for Legifrance's Official Law updates. This is at TESTING stage for now. If you want to follow that topic, you can find more info at https://github.com/ldicarlo/legifrss",
		ManagingEditor: "luca@di-carlo.fr (Luca Di Carlo)",
		LastBuildDate:  now,
		PubDate:        now,
	}

	for _, jorfContent := range result {
		for _, item := range jorfContent.Items {

			for _, step := range item.Container.Structure.Contents {
				feed.Items = tranformHierarchyStep(step, []string{}, feed.Items)
			}
		}
	}

	rss, err := feeds.ToXML(feed)
	if err != nil {
		log.Fatal(err)
	}
	return rss
}

func transformSummary(summary dila.Summary, categories []string) *feeds.RssItem {
	return &feeds.RssItem{
		Title:       summary.Id,
		Link:        "https://www.legifrance.gouv.fr/jorf/id/" + summary.Id,
		Description: summary.Title,
		Category:    strings.Join(categories, "/"),
		Author:      summary.Emitter,
	}
}

func tranformHierarchyStep(hs dila.HierarchyStep, categories []string, result []*feeds.RssItem) []*feeds.RssItem {
	for _, item := range hs.Summaries {
		result = append(result, transformSummary(item, append(categories, hs.Title)))
	}
	for _, nextHs := range hs.Step {
		result = tranformHierarchyStep(nextHs, append(categories, hs.Title), result)
	}
	return result
}
