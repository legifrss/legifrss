package dila

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ldicarlo/legifrss/server/models"
	"github.com/ldicarlo/legifrss/server/utils"
)

func FetchJORF(token string) (str string, lastNJo models.LastNJo) {
	var nbElements = "1"
	fmt.Printf("Fetching the last %s\n", nbElements)
	req, err := http.NewRequest("POST", "https://api.aife.economie.gouv.fr/dila/legifrance-beta/lf-engine-app/consult/lastNJo", strings.NewReader("{\"nbElement\":"+nbElements+"}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	fmt.Println(req)
	if err != nil {
		return "Error while POST request", lastNJo
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return "Error doing request", lastNJo
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error reading body", lastNJo
	}

	json.Unmarshal(body, &lastNJo)

	fmt.Printf("Fetched %d elements\n", len(lastNJo.Containers))
	for _, container := range lastNJo.Containers {
		fmt.Println(container.Id)
	}

	return "", lastNJo
}

func FetchCont(token string, jorfCont string) (str string, joContainerResult models.JOContainerResult) {
	fmt.Printf("Fetching the content for %s\n", jorfCont)
	req, err := http.NewRequest("POST", "https://api.aife.economie.gouv.fr/dila/legifrance-beta/lf-engine-app/consult/jorfCont",
		strings.NewReader("{\"id\":\""+jorfCont+"\",\"pageNumber\":1,\"pageSize\":10}}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	utils.ErrCheck(err)

	resp, err := http.DefaultClient.Do(req)
	utils.ErrCheck(err)

	body, err := ioutil.ReadAll(resp.Body)
	utils.ErrCheck(err)

	json.Unmarshal(body, &joContainerResult)
	fmt.Println(joContainerResult.Items[0].Container.Id + " fetched")
	return "", joContainerResult
}
