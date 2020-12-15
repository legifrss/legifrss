package utils

import (
	"strconv"
	"testing"

	"github.com/ldicarlo/legifrss/server/models"
	"github.com/stretchr/testify/assert"
)

func getArticlesList(startingAt int) []models.JorfArticle {
	return []models.JorfArticle{
		{Content: "2", Order: strconv.Itoa(startingAt + 2)},
		{Content: "1", Order: strconv.Itoa(startingAt + 1)},
		{Content: "3", Order: strconv.Itoa(startingAt + 3)},
		{Content: "4", Order: strconv.Itoa(startingAt + 4)},
	}
}

func TestSortArticles(t *testing.T) {
	unordered := getArticlesList(0)
	sortArticles(unordered)
	assert.Equal(t, "1", unordered[0].Content)
	assert.Equal(t, "2", unordered[1].Content)
	assert.Equal(t, "3", unordered[2].Content)
	assert.Equal(t, "4", unordered[3].Content)

}

func TestSortContent(t *testing.T) {
	input := []models.JorfContainerSection{
		{
			Title:    "3",
			Articles: getArticlesList(10),
			Sections: []models.JorfContainerSection{},
		}, {
			Title:    "2",
			Articles: getArticlesList(5),
			Sections: []models.JorfContainerSection{},
		}, {
			Title:    "4",
			Articles: getArticlesList(15),
			Sections: []models.JorfContainerSection{},
		}, {
			Title:    "1",
			Articles: []models.JorfArticle{},
			Sections: []models.JorfContainerSection{},
		},
	}

	sortContent(input)

	assert.Equal(t, "1", input[0].Title)
	assert.Equal(t, "2", input[1].Title)
	assert.Equal(t, "3", input[2].Title)

}
