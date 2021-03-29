package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/ldicarlo/legifrss/server/pkg/bot"
	"github.com/ldicarlo/legifrss/server/pkg/db"
	"github.com/ldicarlo/legifrss/server/pkg/dila"
	"github.com/ldicarlo/legifrss/server/pkg/models"
	"github.com/ldicarlo/legifrss/server/pkg/token"
	"github.com/ldicarlo/legifrss/server/pkg/utils"
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

func Start() int {
	err, token := token.GetToken(clientId, clientSecret)
	utils.ErrCheckStr(err)

	res := dila.FetchJORF(token)

	var twitterState map[string]models.TwitterJORF

	var jorfContents map[string]models.JOContainerResult
	for _, jorf := range res.Containers {
		nextContent := dila.FetchCont(token, jorf.ID)
		twitterState[jorf.ID] = models.TwitterJORF{
			JORFID:    jorf.ID,
			JORFTitle: jorf.Title,
			URI:       jorf.IDEli,
			StatusID:  0,
			Date:      jorf.Date,
		}
		jorfContents[jorf.ID] = nextContent
	}
	list := utils.ExtractAndConvertDILA(jorfContents)

	total := len(list)
	for i, element := range list {
		fmt.Printf("Fetching the jorf content for %s (%5d/%5d)\n", element.ID, i+1, total)
		result := dila.FetchJorfContent(token, element.ID)
		list[i].Content = utils.ExtractContent(result.Articles, result.Sections)
		twitterState[element]
	}

	db.Persist(list)

	bot.ProcessElems()

	return 0
}

func main() {
	os.Exit(Start())
}
