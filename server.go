package main

import (
	"fmt"

	"github.com/joho/godotenv"
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
}
func Start() string {
	err, token := token.GetToken(clientId, clientSecret)
	if err != "" {
		fmt.Println(err)
	}
	return token
}

func main() {
	fmt.Println(Start())
}
