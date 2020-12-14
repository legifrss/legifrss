package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/ldicarlo/legifrss/server/dila"
	"github.com/ldicarlo/legifrss/server/generate"
	"github.com/ldicarlo/legifrss/server/models"
	"github.com/ldicarlo/legifrss/server/token"
	"github.com/ldicarlo/legifrss/server/utils"
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
	utils.ErrCheckStr(err)

	res := dila.FetchJORF(token)

	var jorfContents []models.JOContainerResult
	for _, jorf := range res.Containers {
		nextContent := dila.FetchCont(token, jorf.Id)
		jorfContents = append(jorfContents, nextContent)
	}
	list := utils.ExtractAndConvertDILA(jorfContents)

	total := len(list)
	for i, element := range list {
		fmt.Printf("Fetching the jorf content for %s (%d/%d)\n", element.Id, i+1, total)
		result := dila.FetchJorfContent(token, element.Id)
		list[i].Content = utils.ExtractContent(result.Articles, result.Sections)

	}
	generate.Generate(list)
	return "", "ok"
}

func main() {
	fmt.Println(Start())
}
