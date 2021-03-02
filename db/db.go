package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ldicarlo/legifrss/server/models"
	"github.com/ldicarlo/legifrss/server/utils"
)

// getAll returns all contents of db.xml
func getAll() (result map[string]models.LegifranceElement) {
	file, err := os.Open("feed/db.json")
	if err != nil {
		return map[string]models.LegifranceElement{}
	}

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &result)
	return result
}

func Query(queryContext models.QueryContext) []models.LegifranceElement {
	var feed = getAll()
	var entries []models.LegifranceElement

	for _, element := range feed {
		if keep(queryContext, &element) {
			entries = append(entries, element)
		}
	}

	return entries
}

func keep(queryContext models.QueryContext, element *models.LegifranceElement) bool {
	if queryContext.Author != "" && !strings.Contains(strings.ToUpper(element.Author), queryContext.Author) {
		return false
	}
	if queryContext.Nature != "" && !strings.Contains(strings.ToUpper(element.Nature), queryContext.Nature) {
		return false
	}
	return strings.Contains(strings.ToUpper(element.Content), queryContext.Keyword) ||
		strings.Contains(strings.ToUpper(element.Title), queryContext.Keyword) ||
		strings.Contains(strings.ToUpper(element.Description), queryContext.Keyword)

}

func Persist(result []models.LegifranceElement) {
	var db = getAll()
	// 864000000000000 nanos = 10 days
	// 86400000000000 nanos = 1 day
	var limitDate = time.Now().Add(-time.Duration(864000000000000))
	var filteredDb = map[string]models.LegifranceElement{}
	for _, element := range db {
		if element.Date.After(limitDate) {
			filteredDb[element.ID] = element
		}
	}

	for _, element := range filteredDb {
		db[element.ID] = element
	}

	jsonResult, _ := json.Marshal(db)
	err := ioutil.WriteFile("db/db.json", jsonResult, 0644)
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
