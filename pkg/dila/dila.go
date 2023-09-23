package dila

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"legifrss/pkg/models"
	"legifrss/pkg/utils"
)

func FetchJORF(token string) (lastNJo models.LastNJo) {
	nbElements := 5
	fmt.Printf("Fetching the last %d\n", nbElements)
	req, err := http.NewRequest("POST", "https://api.piste.gouv.fr/dila/legifrance/lf-engine-app/consult/lastNJo", strings.NewReader("{\"nbElement\":"+strconv.Itoa(nbElements)+"}"))
	utils.ErrCheck(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	utils.ErrCheck(err)
	body, err := ioutil.ReadAll(resp.Body)
	utils.ErrCheck(err)

	json.Unmarshal(body, &lastNJo)

	fmt.Printf("Fetched %d elements\n", len(lastNJo.Containers))

	return lastNJo
}

func FetchCont(token string, jorfCont string) (joContainerResult models.JOContainerResult) {
	fmt.Printf("Fetching the content for %s\n", jorfCont)
	req, err := http.NewRequest("POST", "https://api.piste.gouv.fr/dila/legifrance/lf-engine-app/consult/jorfCont",
		strings.NewReader("{\"id\":\""+jorfCont+"\",\"pageNumber\":1,\"pageSize\":10}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	utils.ErrCheck(err)
	resp, err := http.DefaultClient.Do(req)
	utils.ErrCheck(err)
	err = json.NewDecoder(resp.Body).Decode(&joContainerResult)
	utils.ErrCheck(err)

	return joContainerResult
}

func FetchJorfContent(token string, jorfText string) (joContainerResult models.JorfContainerResult) {
	req, err := http.NewRequest("POST", "https://api.piste.gouv.fr/dila/legifrance/lf-engine-app/consult/jorf",
		strings.NewReader("{\"textCid\":\""+jorfText+"\"}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	utils.ErrCheck(err)

	resp, err := http.DefaultClient.Do(req)
	utils.ErrCheck(err)

	body, err := ioutil.ReadAll(resp.Body)
	utils.ErrCheck(err)

	json.Unmarshal(body, &joContainerResult)
	return joContainerResult
}
