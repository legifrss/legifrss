package utils

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/ldicarlo/legifrss/server/models"
)

func ExtractAndConvertDILA(input []models.JOContainerResult) (result []models.LegifranceElement) {
	loc, err := time.LoadLocation("Europe/Paris")
	ErrCheck(err)
	for _, jorfContent := range input {
		for _, item := range jorfContent.Items {
			ts := time.Unix(item.Container.Timestamp/1000, 0).In(loc)
			for _, step := range item.Container.Structure.Contents {
				result = tranformHierarchyStep(step, []string{}, result, ts)
			}
		}
	}
	return
}

func transformSummary(summary models.Summary, categories []string, publicationDate time.Time) models.LegifranceElement {
	return models.LegifranceElement{
		Id:          summary.Id,
		Title:       summary.Id,
		Link:        "https://www.legifrance.gouv.fr/jorf/id/" + summary.Id,
		Description: summary.Title,
		Category:    categories,
		Author:      summary.Emitter,
		Nature:      summary.Nature,
		Date:        publicationDate.String(),
	}
}

func tranformHierarchyStep(
	hs models.HierarchyStep,
	categories []string,
	result []models.LegifranceElement,
	pub time.Time,
) []models.LegifranceElement {

	for _, item := range hs.Summaries {
		result = append(result, transformSummary(item, append(categories, hs.Title), pub))
	}
	for _, nextHs := range hs.Step {
		result = tranformHierarchyStep(nextHs, append(categories, hs.Title), result, pub)
	}
	return result
}

func ErrCheck(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic("Panic.")
	}
}

func ErrCheckStr(str string) {
	if str != "" {
		fmt.Println(str)
		panic("Panic.")
	}
}

func ExtractContent(articles []models.JorfArticle, sections []models.JorfContainerSection) string {
	str := ""
	sortArticles(articles)
	for _, article := range articles {
		str += "Article " + article.Order + "<br />" + article.Content
	}
	sortContent(sections)
	for _, section := range sections {
		str += section.Title
		str += ExtractContent(section.Articles, section.Sections)

	}
	return str
}

func sortContent(sections []models.JorfContainerSection) {

	sort.Sort(SortByArticleOrder(sections))
}

func sortArticles(articles []models.JorfArticle) {
	sort.Sort(SortByOrder(articles))
}

func toInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return -1
	}
	return i
}

type SortByOrder []models.JorfArticle

func (a SortByOrder) Len() int           { return len(a) }
func (a SortByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByOrder) Less(i, j int) bool { return toInt(a[i].Order) < toInt(a[j].Order) }

type SortByArticleOrder []models.JorfContainerSection

func (a SortByArticleOrder) Len() int      { return len(a) }
func (a SortByArticleOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByArticleOrder) Less(i, j int) bool {

	return lowArticleOrderInSection(a[i]) < lowArticleOrderInSection(a[j])
}

func lowArticleOrderInSection(section models.JorfContainerSection) int {
	if len(section.Articles) == 0 {
		return -1
	}
	return toInt(section.Articles[0].Order)

}
