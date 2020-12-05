package dila

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func FetchJORF(token string) (str string, lastNJo LastNJo) {
	req, err := http.NewRequest("POST", "https://api.aife.economie.gouv.fr/dila/legifrance-beta/lf-engine-app/consult/lastNJo", strings.NewReader("{\"nbElement\":5}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
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
	return "", lastNJo
}

func FetchCont(token string, jorfCont string) (str string, joContainerResult JOContainerResult) {
	req, err := http.NewRequest("POST", "https://api.aife.economie.gouv.fr/dila/legifrance-beta/lf-engine-app/consult/jorfCont",
		strings.NewReader("{\"textCid\":\""+jorfCont+"\"}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return "Error while POST request", joContainerResult
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return "Error doing request", joContainerResult
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return "Error reading body", joContainerResult
	}
	json.Unmarshal(body, &joContainerResult)
	return "", joContainerResult
}
