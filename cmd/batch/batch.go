package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"legifrss/pkg/db"
	"legifrss/pkg/dila"
	"legifrss/pkg/models"
	"legifrss/pkg/token"
	"legifrss/pkg/utils"
)

var clientId string
var clientSecret string

func init() {
	val, ok := os.LookupEnv("ENV_FILE")
	if !ok  {
		val = ".env"
	}
	envs, err := godotenv.Read(val)
	if err != nil {
		panic(fmt.Sprintf("missing env file, reading from %s, %s", val, err))
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

	elementsWithContents := map[string]models.JORFElement{}
	i := 0
	for jorfID, jorf := range elements {
		j := 0
		for jorfContentID, content := range jorf.JORFContents {
			j++
			i++
			fmt.Printf("Fetching the jorf content for %s (%5d/%5d)\n", content.ID, i+1, total)
			result := dila.FetchJorfContent(token, content.ID)
			content.Content = utils.ExtractContent(result.Articles, result.Sections)
			jorf.JORFContents[jorfContentID] = content
		}
		elementsWithContents[jorfID] = jorf
	}

	db.Persist(elementsWithContents)
	db.PersistTwitterState(twitterState)
	// bot.ProcessElems()
	// bot.CleanOldElems()
	return 0
}

func main() {
	os.Exit(Start())
}
