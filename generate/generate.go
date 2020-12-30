package generate

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/feeds"
	"github.com/ldicarlo/legifrss/server/models"
	"github.com/ldicarlo/legifrss/server/rss"
	"github.com/ldicarlo/legifrss/server/utils"
)

func Generate(result []models.LegifranceElement) {
	for i, element := range result {
		result[i].SanitizedAuthor = sanitizeName(element.Author)
		result[i].SanitizedNature = sanitizeName(element.Nature)
	}

	natureMap := mapByNature(result)
	for key, posts := range natureMap {
		filePath := key + "_all.xml"
		feed := rss.TransformToRSS(posts, models.FeedDescription{TitleSuffix: "- " + key, LinkSuffix: filePath, DescriptionSuffix: ""})
		appendToFile(feed, filePath)

	}
	authorMap := mapByAuthor(result)
	for key, posts := range authorMap {
		filePath := "all_" + key + ".xml"
		feed := rss.TransformToRSS(posts, models.FeedDescription{TitleSuffix: "- " + key, LinkSuffix: filePath, DescriptionSuffix: ""})

		appendToFile(feed, filePath)

	}
	authorNatureMap := mapByAuthorAndNature(result)
	for key1, maps := range authorNatureMap {

		for key2, item := range maps {
			filePath := key2 + "_" + key1 + ".xml"
			feed := rss.TransformToRSS(item, models.FeedDescription{TitleSuffix: "- " + key1 + " - " + key2, LinkSuffix: filePath, DescriptionSuffix: ""})
			appendToFile(feed, filePath)
		}
	}

	globalFeed := rss.TransformToRSS(result, models.FeedDescription{TitleSuffix: "", LinkSuffix: "all.xml", DescriptionSuffix: ""})
	f, err := os.Create("feed/all.xml")
	utils.ErrCheck(err)

	elem, err := feeds.ToXML(globalFeed)
	utils.ErrCheck(err)

	f.WriteString(elem)

}

func mapByNature(input []models.LegifranceElement) map[string][]models.LegifranceElement {
	var result = make(map[string][]models.LegifranceElement)
	for _, item := range input {
		if result[item.SanitizedNature] != nil {
			result[item.SanitizedNature] = append(result[item.SanitizedNature], item)
		} else {
			result[item.SanitizedNature] = []models.LegifranceElement{item}
		}
	}
	return result
}

func mapByAuthor(input []models.LegifranceElement) map[string][]models.LegifranceElement {
	var result = make(map[string][]models.LegifranceElement)
	for _, item := range input {
		if result[item.SanitizedAuthor] != nil {
			result[item.SanitizedAuthor] = append(result[item.SanitizedAuthor], item)
		} else {
			result[item.SanitizedAuthor] = []models.LegifranceElement{item}
		}
	}
	return result
}

func sanitizeName(str string) string {
	if str == "" {
		return "unknown"
	}
	str = strings.ToLower(str)
	replacer := strings.NewReplacer("  ", " ", ",", "", " ", "-", "'", "-", "--", "-", "è", "e", "é", "e", "ç", "c", "à", "a", "ô", "o", "û", "u")
	newStr := str
	for newStr != replacer.Replace(newStr) {
		newStr = replacer.Replace(newStr)
	}
	return customReplacements(newStr)

}

func customReplacements(str string) string {
	replacements := map[string]string{
		"ministere-de-l-education-nationale-de-la-jeunesse-et-des-sports-sports.xml": "ministere-de-l-education-nationale-de-la-jeunesse-et-des-sports.xml",
		"ministere-du-travail-de-l-emploi-et-de-l-insertion-insertion.xml":           "ministere-du-travail-de-l-emploi-et-de-l-insertion.xml",
	}
	if value, found := replacements[str]; found {
		return value
	}
	return str
}

func mapByAuthorAndNature(input []models.LegifranceElement) map[string]map[string][]models.LegifranceElement {
	var result = make(map[string]map[string][]models.LegifranceElement)
	for _, item := range input {
		if result[item.SanitizedAuthor] == nil {
			result[item.SanitizedAuthor] = make(map[string][]models.LegifranceElement)
			result[item.SanitizedAuthor][item.SanitizedNature] = []models.LegifranceElement{item}

		} else {
			if result[item.SanitizedAuthor][item.SanitizedNature] == nil {
				result[item.SanitizedAuthor][item.SanitizedNature] = []models.LegifranceElement{item}
			} else {
				result[item.SanitizedAuthor][item.SanitizedNature] = append(result[item.SanitizedAuthor][item.SanitizedNature], item)
			}

		}
	}

	return result
}

func appendToFile(feed *feeds.AtomFeed, path string) {
	path = "feed/" + path
	limit := 50
	fmt.Println("path to print " + path)
	fmt.Println("entries " + strconv.Itoa(len(feed.Entries)))
	if _, err := os.Stat(path); err != nil || len(feed.Entries) > limit {
		write(feed, path)
	}

	file, err := os.Open(path)
	utils.ErrCheck(err)

	byteValue, _ := ioutil.ReadAll(file)

	var oldFeed feeds.AtomFeed
	xml.Unmarshal(byteValue, &oldFeed)
	fmt.Println("Merging feed : " + path + ",entries: " + strconv.Itoa(len(feed.Entries)))

	fmt.Println(oldFeed.Entries)

	feed = mergeFeeds(feed, oldFeed, limit)
	write(feed, path)
}

func mergeFeeds(newFeed *feeds.AtomFeed, oldFeed feeds.AtomFeed, limit int) *feeds.AtomFeed {
	list := append(newFeed.Entries, oldFeed.Entries...)
	if len(list) > limit {
		list = list[:limit]
	}
	newFeed.Entries = list
	return newFeed
}

func write(feed *feeds.AtomFeed, path string) {
	f, err := os.Create(path)
	utils.ErrCheck(err)

	elem, err := feeds.ToXML(feed)
	utils.ErrCheck(err)

	f.WriteString(elem)
}
