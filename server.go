package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/ldicarlo/legifrss/server/dila"
	"github.com/ldicarlo/legifrss/server/token"
)

var clientId string
var clientSecret string

func init() {
	envs, err := godotenv.Read(".env")
	if err != nil {
		panic("missing env file")
	}
	clientId = envs["client_id"]
	clientSecret = envs["client_secret"]

	if clientId == "" || clientSecret == "" {
		panic("Missing one of the env params")
	}

}

func Start() (str string, result string) {
	err, token := token.GetToken(clientId, clientSecret)
	if err != "" {
		return err, ""
	}

	err, res := dila.FetchJORF(token)
	if err != "" {
		return err, ""
	}
	err, res2 := dila.FetchCont(token, res.Containers[0].Id)
	if err != "" {
		return err, ""
	}
	fmt.Println(res2)
	return "", "ok"
}

func main() {
	fmt.Println(Start())
}
