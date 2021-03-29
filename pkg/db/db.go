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
	return
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

func ExtractContentToPublish() []models.TwitterJORF {
	var twitterState = FetchTwitterStates()
	var entries []models.TwitterJORF

	for _, element := range twitterState {
		if hasAnyMissing(element) {
			entries = append(entries, element)
		}
	}

	return entries
}

func hasAnyMissing(jorf models.TwitterJORF) bool {
	if jorf.StatusID == 0 {
		return true
	}
	for _, element := range jorf.JORFContents {
		if element.StatusID == 0 {
			return true
		}
	}
	return false
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

func PersistTwitterState(jorfs []models.TwitterJORF) {
	existing := FetchTwitterStates()
	var limitDate = time.Now().Add(-time.Duration(865000000000000))
	var filteredDb = map[string]models.TwitterJORF{}
	for _, element := range existing {
		if element.Date.After(limitDate) {
			filteredDb[element.JORFID] = element
		}
	}

	for _, element := range jorfs {
		if val, ok := filteredDb[element.JORFID]; ok {
			filteredDb[element.JORFID] = mergeTwitterJORFs(val, element)
		} else {
			filteredDb[element.JORFID] = element
		}
	}

	jsonResult, _ := json.Marshal(filteredDb)
	err := ioutil.WriteFile("db/twitter_states.json", jsonResult, 0644)
	utils.ErrCheck(err)
}

func FetchTwitterStates() (jorfs map[string]models.TwitterJORF) {
	file, err := os.Open("db/twitter_states.json")
	if err != nil {
		return map[string]models.TwitterJORF{}
	}

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &jorfs)
	return
}

func mergeTwitterJORFs(before models.TwitterJORF, after models.TwitterJORF) (result models.TwitterJORF) {
	result.Date = after.Date
	result.JORFID = after.JORFID
	result.JORFTitle = after.JORFID
	result.StatusID = max(before.StatusID, after.StatusID)
	result.JORFTitle = after.JORFTitle
	result.URI = after.URI
	contents := before.JORFContents
	for _, content := range after.JORFContents {
		if val, ok := contents[content.JORFContentID]; ok {
			contents[content.JORFContentID] = mergeTwitterJORFContents(val, content)
		} else {
			contents[content.JORFContentID] = content
		}
	}
	return
}

func mergeTwitterJORFContents(before models.TwitterJORFContent, after models.TwitterJORFContent) (result models.TwitterJORFContent) {
	result.Content = after.Content
	result.JORFContentID = after.JORFContentID
	result.StatusID = max(before.StatusID, after.StatusID)
	return
}

func max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}
