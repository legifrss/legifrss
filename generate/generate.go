package generate

import (
	"os"
	"strings"

	"github.com/gorilla/feeds"
	"github.com/ldicarlo/legifrss/server/models"
	"github.com/ldicarlo/legifrss/server/rss"
	"github.com/ldicarlo/legifrss/server/utils"
)

func Generate(result []models.LegifranceElement) {
	for i, element := range result {
		result[i].Author = sanitizeName(element.Author)
		result[i].Nature = sanitizeName(element.Nature)
	}

	natureMap := mapByNature(result)
	for key, posts := range natureMap {

		feed := rss.TransformToRSS(posts, models.FeedDescription{TitleSuffix: "- " + key, LinkSuffix: "", DescriptionSuffix: ""})
		f, err := os.Create("feed/" + key + "_all.xml")

		utils.ErrCheck(err)

		elem, err := feeds.ToXML(feed)
		utils.ErrCheck(err)

		f.WriteString(elem)

	}
	authorMap := mapByAuthor(result)
	for key, posts := range authorMap {
		feed := rss.TransformToRSS(posts, models.FeedDescription{TitleSuffix: "- " + key, LinkSuffix: "", DescriptionSuffix: ""})

		f, err := os.Create("feed/all_" + key + ".xml")
		utils.ErrCheck(err)

		elem, err := feeds.ToXML(feed)
		utils.ErrCheck(err)

		f.WriteString(elem)

	}
	authorNatureMap := mapByAuthorAndNature(result)
	for key1, maps := range authorNatureMap {

		for key2, item := range maps {

			feed := rss.TransformToRSS(item, models.FeedDescription{TitleSuffix: "- " + key1 + " - " + key2, LinkSuffix: "", DescriptionSuffix: ""})
			f, err := os.Create("feed/" + key2 + "_" + key1 + ".xml")
			utils.ErrCheck(err)

			elem, err := feeds.ToXML(feed)
			utils.ErrCheck(err)

			f.WriteString(elem)
		}
	}

	globalFeed := rss.TransformToRSS(result, models.FeedDescription{TitleSuffix: "", LinkSuffix: "", DescriptionSuffix: ""})
	f, err := os.Create("feed/all.xml")
	utils.ErrCheck(err)

	elem, err := feeds.ToXML(globalFeed)
	utils.ErrCheck(err)

	f.WriteString(elem)

}

func mapByNature(input []models.LegifranceElement) map[string][]models.LegifranceElement {
	var result = make(map[string][]models.LegifranceElement)
	for _, item := range input {
		if result[item.Nature] != nil {
			result[item.Nature] = append(result[item.Nature], item)
		} else {
			result[item.Nature] = []models.LegifranceElement{item}
		}
	}
	return result
}

func mapByAuthor(input []models.LegifranceElement) map[string][]models.LegifranceElement {
	var result = make(map[string][]models.LegifranceElement)
	for _, item := range input {
		if result[item.Author] != nil {
			result[item.Author] = append(result[item.Author], item)
		} else {
			result[item.Author] = []models.LegifranceElement{item}
		}
	}
	return result
}

func sanitizeName(str string) string {
	if str == "" {
		return "unknown"
	}
	str = strings.ToLower(str)
	replacer := strings.NewReplacer("  ", " ", ",", "", " ", "-", "'", "-", "--", "-")
	newStr := str
	for newStr != replacer.Replace(newStr) {
		newStr = replacer.Replace(newStr)

	}
	return customReplacements(newStr)

}

func customReplacements(str string) string {
	replacements := map[string]string{
		"ministère-de-l-éducation-nationale-de-la-jeunesse-et-des-sports-sports.xml": "ministère-de-l-éducation-nationale-de-la-jeunesse-et-des-sports.xml",
		"ministère-du-travail-de-l-emploi-et-de-l-insertion-insertion.xml":           "ministère-du-travail-de-l-emploi-et-de-l-insertion.xml",
	}
	if value, found := replacements[str]; found {
		return value
	}
	return str
}

func mapByAuthorAndNature(input []models.LegifranceElement) map[string]map[string][]models.LegifranceElement {
	var result = make(map[string]map[string][]models.LegifranceElement)
	for _, item := range input {
		if result[item.Author] == nil {
			result[item.Author] = make(map[string][]models.LegifranceElement)
			result[item.Author][item.Nature] = []models.LegifranceElement{item}

		} else {
			if result[item.Author][item.Nature] == nil {
				result[item.Author][item.Nature] = []models.LegifranceElement{item}
			} else {
				result[item.Author][item.Nature] = append(result[item.Author][item.Nature], item)
			}

		}
	}

	return result
}

func appendToFile(feed feeds.AtomFeed, path string) {
	if _, err := os.Stat(path); err == nil || len(feed.Entries) > 500 {
		write(feed, path)
	}

	// loadExisting()

	// append oldfeed.Items to feed.Items
	// keep limit
}

func write(feed feeds.AtomFeed, path string) {
	f, err := os.Create(path)
	utils.ErrCheck(err)

	elem, err := feeds.ToXML(&feed)
	utils.ErrCheck(err)

	f.WriteString(elem)
}
