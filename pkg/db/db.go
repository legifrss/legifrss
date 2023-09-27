package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
	"legifrss/pkg/models"
	"legifrss/pkg/utils"
)

func PersistToken(token oauth1.Token) {
	jsonResult, _ := json.Marshal(token)
	err := ioutil.WriteFile("token.json", jsonResult, 0644)
	utils.ErrCheck(err)
}

func GetToken() (token oauth1.Token, noFile error) {
	file, err := os.Open("token.json")
	if err != nil {
		noFile = err
		return
	}
	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &token)
	return token, nil

}

func getAll() (result map[string]models.JORFElement) {
	file, err := os.Open("/home/legifrss/db.json")
	if err != nil {
		return map[string]models.JORFElement{}
	}

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &result)
	return
}

func getAllLegifranceElements() (result []models.LegifranceElement) {
	db := getAll()
	for _, jorfs := range db {
		for _, content := range jorfs.JORFContents {
			result = append(result, content)
		}
	}
	return
}

func Query(queryContext models.QueryContext) []models.LegifranceElement {
	var feed = getAll()
	var entries []models.LegifranceElement

	for _, jorfElement := range feed {
		for _, element := range jorfElement.JORFContents {
			if keep(queryContext, &element) {
				entries = append(entries, element)
			}
		}
	}

	return entries
}

func fetchAllJORFKeys(keys []string) map[string]models.JORFElement {
	db := getAll()
	toPublish := map[string]models.JORFElement{}
	for key, elem := range db {
		if contains(keys, key) {
			toPublish[key] = elem
		}
	}
	return toPublish
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ExtractContentToPublish() (map[string]models.JORFElement, map[string]models.TwitterJORF) {
	var twitterState = FetchTwitterStates()
	state := map[string]models.TwitterJORF{}
	jorfKeysToFetch := []string{}
	for ID, element := range twitterState {
		if hasAnyMissing(element) {
			state[ID] = element
			jorfKeysToFetch = append(jorfKeysToFetch, ID)
		}
	}

	toPublish := fetchAllJORFKeys(jorfKeysToFetch)

	return toPublish, state
}

func hasAnyMissing(jorf models.TwitterJORF) bool {
	if jorf.StatusID == 0 {
		return true
	}
	for _, element := range jorf.JORFContents {
		if element == 0 {
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

func Persist(result map[string]models.JORFElement) {
	var db = getAll()
	// 864000000000000 nanos = 10 days
	// 86400000000000 nanos = 1 day
	var limitDate = time.Now().Add(-time.Duration(864000000000000))
	var filteredDb = map[string]models.JORFElement{}
	for _, element := range db {
		if element.Date.After(limitDate) {
			filteredDb[element.JORFID] = element
		}
	}

	for _, element := range result {
		filteredDb[element.JORFID] = element
	}

	jsonResult, _ := json.Marshal(filteredDb)
	err := ioutil.WriteFile("/home/legifrss/db.json", jsonResult, 0644)
	utils.ErrCheck(err)
}

func GetAuthors() (arr []string) {
	values := getAllLegifranceElements()
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
	values := getAllLegifranceElements()
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

func PersistTwitterState(jorfs map[string]models.TwitterJORF) {
	existing := FetchTwitterStates()
	filteredDb := map[string]models.TwitterJORF{}
	for ID, element := range existing {
		filteredDb[ID] = element
	}

	for ID, element := range jorfs {
		if val, ok := filteredDb[ID]; ok {
			filteredDb[ID] = mergeTwitterJORFs(val, element)
		} else {
			filteredDb[ID] = element
		}
	}

	jsonResult, _ := json.Marshal(filteredDb)
	err := ioutil.WriteFile("twitter_states.json", jsonResult, 0644)
	utils.ErrCheck(err)
}

func OverrideTwitterStates(jorfs map[string]models.TwitterJORF) {
	jsonResult, _ := json.Marshal(jorfs)
	err := ioutil.WriteFile("twitter_states.json", jsonResult, 0644)
	utils.ErrCheck(err)
}

func FetchTwitterStates() (jorfs map[string]models.TwitterJORF) {
	file, err := os.Open("twitter_states.json")
	if err != nil {
		return map[string]models.TwitterJORF{}
	}

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &jorfs)
	return
}

func mergeTwitterJORFs(before models.TwitterJORF, after models.TwitterJORF) (result models.TwitterJORF) {
	result.StatusID = max(before.StatusID, after.StatusID)
	contents := before.JORFContents
	for ID, content := range after.JORFContents {
		if val, ok := contents[ID]; ok {
			contents[ID] = mergeTwitterJORFContents(val, content)
		} else {
			contents[ID] = content
		}
	}
	result.JORFContents = contents
	return
}

func mergeTwitterJORFContents(before int64, after int64) int64 {
	return max(before, after)

}

func max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

func GetAllKeys() (keys []string) {
	db := getAll()
	for k := range db {
		keys = append(keys, k)
	}
	return

}
