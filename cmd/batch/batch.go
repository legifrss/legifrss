package main

import (
	"fmt"
	"os"
	"time"

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

	elements := map[string]models.JORFElement{}
	total := 0
	loc, locErr := time.LoadLocation("Europe/Paris")
	utils.ErrCheck(locErr)

	twitterState := map[string]models.TwitterJORF{}

	for _, jorf := range res.Containers {
		nextContent := dila.FetchCont(token, jorf.ID)
		elements[jorf.ID] = utils.ExtractAndConvertDILA(nextContent, models.JORFElement{
			JORFID:    jorf.ID,
			JORFTitle: jorf.Title,
			Date:      time.Unix(jorf.Date/1000, 0).In(loc),
			URI:       jorf.IDEli,
		})
		keys := map[string]int64{}
		for _, elem := range elements[jorf.ID].JORFContents {
			keys[elem.ID] = 0
		}
		twitterState[jorf.ID] = models.TwitterJORF{
			StatusID:     0,
			JORFContents: keys,
		}
		total = total + len(elements[jorf.ID].JORFContents)
	}
	i := 0
	for _, jorf := range elements {
		j := 0
		for _, content := range jorf.JORFContents {
			j++
			i++
			fmt.Printf("Fetching the jorf content for %s (%5d/%5d)\n", content.ID, i+1, total)
			result := dila.FetchJorfContent(token, content.ID)
			content.Content = utils.ExtractContent(result.Articles, result.Sections)

		}
	}

	db.Persist(elements)
	db.PersistTwitterState(twitterState)
	bot.ProcessElems()

	return 0
}

func main() {
	os.Exit(Start())
}
