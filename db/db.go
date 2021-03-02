package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ldicarlo/legifrss/server/models"
	"github.com/ldicarlo/legifrss/server/utils"
)

// GetAll returns all contents of db.xml
func GetAll() (result []models.LegifranceElement) {
	file, err := os.Open("feed/db.json")
	utils.ErrCheck(err)

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &result)
	return result
}

func Query(queryContext models.QueryContext) []models.LegifranceElement {
	var feed = GetAll()
	var entries []models.LegifranceElement

	for _, element := range feed {
		if keep(queryContext, &element) {
			entries = append(entries, element)
		}
	}

	return entries
}

func keep(queryContext models.QueryContext, element *models.LegifranceElement) bool {
	if queryContext.Author != "" && !strings.Contains(element.Author, queryContext.Author) {
		return false
	}
	if queryContext.Nature != "" && !strings.Contains(element.Nature, queryContext.Nature) {
		return false
	}
	return strings.Contains(element.Content, queryContext.Keyword) || strings.Contains(element.Title, queryContext.Keyword) || strings.Contains(element.Description, queryContext.Keyword)

}

func Persist(result []models.LegifranceElement) {
	merged := merge(result,
		// GetAll(),
		[]models.LegifranceElement{},
		10000)
	jsonResult, _ := json.Marshal(merged)
	err := ioutil.WriteFile("feed/db.json", jsonResult, 0644)
	utils.ErrCheck(err)
}

func merge(newElements []models.LegifranceElement, oldElements []models.LegifranceElement, limit int) (mergeResult []models.LegifranceElement) {
	list := append(newElements, oldElements...)
	if len(list) > limit {
		list = list[:limit]
	}
	mergeResult = list
	return mergeResult
}
