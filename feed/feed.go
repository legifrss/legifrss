package feed

import (
	"fmt"
	"os"

	"github.com/gorilla/feeds"
	"github.com/ldicarlo/legifrss/server/models"
	"github.com/ldicarlo/legifrss/server/rss"
	"github.com/ldicarlo/legifrss/server/utils"
)

func Generate(result []models.LegifranceElement) {
	natureMap := mapByNature(result)
	//	authorMap := mapByAuthor(result)
	for k, posts := range natureMap {
		//fmt.Println(posts)
		feed := rss.TransformToRSS(posts, models.FeedDescription{TitleSuffix: "", LinkSuffix: "", DescriptionSuffix: ""})
		f, err := os.Create("feed/nature-" + k + ".xml")

		utils.ErrCheck(err)

		elem, err := feeds.ToXML(feed)
		utils.ErrCheck(err)

		f.WriteString(elem)

	}
	authorMap := mapByAuthor(result)
	//	authorMap := mapByAuthor(result)
	for k, posts := range authorMap {
		//fmt.Println(posts)
		feed := rss.TransformToRSS(posts, models.FeedDescription{TitleSuffix: "", LinkSuffix: "", DescriptionSuffix: ""})
		f, err := os.Create("feed/author-" + k + ".xml")

		utils.ErrCheck(err)

		elem, err := feeds.ToXML(feed)
		utils.ErrCheck(err)

		f.WriteString(elem)

	}

	globalFeed := rss.TransformToRSS(result, models.FeedDescription{TitleSuffix: "", LinkSuffix: "", DescriptionSuffix: ""})
	f, err := os.Create("feed/feed.xml")
	utils.ErrCheck(err)

	elem, err := feeds.ToXML(globalFeed)
	utils.ErrCheck(err)

	f.WriteString(elem)

}

func mapByNature(input []models.LegifranceElement) map[string][]models.LegifranceElement {
	var result = make(map[string][]models.LegifranceElement)
	for _, item := range input {
		fmt.Println(item.Nature)
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
