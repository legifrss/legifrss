package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"legifrss/pkg/dila"
	"legifrss/pkg/token"
	"legifrss/pkg/utils"
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

func main() {
	err, token := token.GetToken(clientId, clientSecret)
	utils.ErrCheckStr(err)
	println(token)
	println(os.Args[1])
	fmt.Println(dila.FetchCont(token, os.Args[1]))
}
