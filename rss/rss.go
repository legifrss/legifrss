package rss

import (
	"log"
	"time"

	"github.com/gorilla/feeds"
	"github.com/ldicarlo/legifrss/server/dila"
)

// TransformToRSS is the Main transformation function
func TransformToRSS(result []dila.JOContainerResult) string {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "Legifrance RSS",
		Link:        &feeds.Link{Href: "https://legifrance.gouv.fr"},
		Description: "This is a non-official RSS feed for Legifrance. This is at TESTING stage for now. If you want to follow that topic, you can find more info at https://github.com/ldicarlo/legifrss",
		Author:      &feeds.Author{Name: "Luca Di Carlo", Email: "luca@di-carlo.fr"},
		Updated:     now,
	}

	for _, jorfContent := range result {
		for _, item := range jorfContent.Items {
			for _, step := range item.Container.Structure.Contents {
				feed.Items = tranformHierarchyStep(step, feed.Items)
			}
		}
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}
	return rss
}

func transformSummary(summary dila.Summary) *feeds.Item {
	return &feeds.Item{
		Title:       summary.Id,
		Link:        &feeds.Link{Href: "https://www.legifrance.gouv.fr/jorf/id/" + summary.Id},
		Description: summary.Title,
	}
}

func tranformHierarchyStep(hs dila.HierarchyStep, result []*feeds.Item) []*feeds.Item {
	for _, item := range hs.Summaries {
		result = append(result, transformSummary(item))
	}
	for _, nextHs := range hs.Step {
		result = tranformHierarchyStep(nextHs, result)
	}
	return result
}
