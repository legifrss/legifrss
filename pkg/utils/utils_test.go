package utils

import (
	"strconv"
	"testing"

	"github.com/ldicarlo/legifrss/server/pkg/models"
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
	// JORFTEXT000042665307 is a good complex example.
	input := []models.JorfContainerSection{
		{
			Title:    "3",
			Articles: getArticlesList(10),
			Sections: []models.JorfContainerSection{},
		}, {
			Title:    "2",
			Articles: []models.JorfArticle{},
			Sections: []models.JorfContainerSection{{
				Title:    "Haha",
				Sections: []models.JorfContainerSection{},
				Articles: getArticlesList(5)}},
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

func TestPrepareTweetContent(t *testing.T) {
	assert.Equal(t, "The great...", PrepareTweetContent("The great revolution", 9))
	assert.Equal(t, "The ", PrepareTweetContent("The ", 9))
	assert.Equal(t, "The greaté...", PrepareTweetContent("The greaté revolution", 10))
	assert.Equal(t, "The great...", PrepareTweetContent("The greaté", 9))
	assert.Equal(t, "The great revol...", PrepareTweetContent("The great revolution", 15))
	assert.Equal(t, "éééé", PrepareTweetContent("éééé", 5))
}

func TestContains(t *testing.T) {
	assert.True(t, contains([]string{"Test", "Test3"}, "Test"))
	assert.False(t, contains([]string{"Test", "Test3"}, "Test2"))
}

func TestCleanNonExistingKeys(t *testing.T) {
	array := map[string]models.TwitterJORF{
		"FIRST": {
			StatusID: 100, JORFContents: map[string]int64{},
		},
		"SECOND": {
			StatusID: 100, JORFContents: map[string]int64{},
		},
		"THIRD": {
			StatusID: 100, JORFContents: map[string]int64{},
		},
	}

	assert.Equal(t, 2, len(CleanNonExistingKeys(array, []string{"FIRST", "SECOND"})))

}
