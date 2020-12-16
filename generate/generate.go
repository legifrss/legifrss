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
	natureMap := mapByNature(result)
	for key, posts := range natureMap {
		var k = sanitizeName(key)

		feed := rss.TransformToRSS(posts, models.FeedDescription{TitleSuffix: "- " + k, LinkSuffix: "", DescriptionSuffix: ""})
		f, err := os.Create("feed/" + k + "_all.xml")

		utils.ErrCheck(err)

		elem, err := feeds.ToXML(feed)
		utils.ErrCheck(err)

		f.WriteString(elem)

	}
	authorMap := mapByAuthor(result)
	for key, posts := range authorMap {
		var k = sanitizeName(key)
		feed := rss.TransformToRSS(posts, models.FeedDescription{TitleSuffix: "- " + k, LinkSuffix: "", DescriptionSuffix: ""})
		f, err := os.Create("feed/all_" + k + ".xml")

		utils.ErrCheck(err)

		elem, err := feeds.ToXML(feed)
		utils.ErrCheck(err)

		f.WriteString(elem)

	}
	authorNatureMap := mapByAuthorAndNature(result)
	for key1, maps := range authorNatureMap {
		k1 := sanitizeName(key1)
		for key2, item := range maps {
			k2 := sanitizeName(key2)
			feed := rss.TransformToRSS(item, models.FeedDescription{TitleSuffix: "- " + k1 + " - " + k2, LinkSuffix: "", DescriptionSuffix: ""})
			f, err := os.Create("feed/" + k2 + "_" + k1 + ".xml")

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
	return newStr

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
