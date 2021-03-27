package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/ldicarlo/legifrss/server/pkg/models"
	"github.com/ldicarlo/legifrss/server/pkg/utils"
)

func PersistToken(token oauth1.Token) {
	jsonResult, _ := json.Marshal(token)
	err := ioutil.WriteFile("db/token.json", jsonResult, 0644)
	utils.ErrCheck(err)
}

func GetToken() (token oauth1.Token, noFile error) {
	file, err := os.Open("db/token.json")
	if err != nil {
		noFile = err
		return
	}
	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &token)
	return token, nil

}

// getAll returns all contents of db.json
func getAll() (result map[string]models.LegifranceElement) {
	file, err := os.Open("db/db.json")
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

func ExtractContentToPublish() []models.LegifranceElement {
	var feed = getAll()
	var entries []models.LegifranceElement

	for _, element := range feed {
		if element.TwitterPublished == 0 {
			entries = append(entries, element)
			if len(entries) >= 99 {
				break
			}
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

	for _, element := range result {
		if val, ok := filteredDb[element.ID]; ok {
			element.TwitterPublished = val.TwitterPublished
		}
		filteredDb[element.ID] = element
	}

	jsonResult, _ := json.Marshal(filteredDb)
	err := ioutil.WriteFile("db/db.json", jsonResult, 0644)
	utils.ErrCheck(err)
}

func GetAuthors() (arr []string) {
	values := getAll()
	authors := map[string]string{}
	for _, value := range values {
		authors[value.Author] = ""
	}

	for key := range authors {
		if key != "" {
			arr = append(arr, key)
		}
	}
	sort.Strings(arr[:])
	return arr
}

func GetNatures() (arr []string) {
	values := getAll()
	natures := map[string]string{}
	for _, value := range values {
		natures[value.Nature] = ""
	}

	for key := range natures {
		if key != "" {
			arr = append(arr, key)
		}
	}
	sort.Strings(arr[:])
	return arr
}
