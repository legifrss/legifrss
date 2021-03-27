package token

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type jwt struct {
	AccessToken string `json:"access_token"`
}

func GetToken(clientId string, clientSecret string) (string, string) {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", clientId)
	form.Add("client_secret", clientSecret)
	form.Add("scope", "openid")

	resp, err := http.Post("https://oauth.aife.economie.gouv.fr/api/oauth/token/api/oauth/token", "application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		return "Error while POST request", ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error reading body", ""
	}
	var jwt jwt
	json.Unmarshal(body, &jwt)
	//fmt.Printf("Secret Token is %s\n", jwt.AccessToken)
	return "", jwt.AccessToken

}
