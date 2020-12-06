package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/ldicarlo/legifrss/server/dila"
	"github.com/ldicarlo/legifrss/server/rss"
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
	var jorfContents []dila.JOContainerResult
	for _, jorf := range res.Containers {
		err, nextContent := dila.FetchCont(token, jorf.Id)
		if err != "" {
			continue
		}
		jorfContents = append(jorfContents, nextContent)
	}

	feed := rss.TransformToRSS(jorfContents)
	f, er := os.Create("feed/feed.xml")
	if er != nil {
		fmt.Println(err)
		return err, ""
	}
	f.WriteString(feed)
	return "", "ok"
}

func main() {
	fmt.Println(Start())
}
