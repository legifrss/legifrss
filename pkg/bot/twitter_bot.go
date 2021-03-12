package bot

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	oauthTwitter "github.com/dghubble/oauth1/twitter"
	"github.com/joho/godotenv"
	"github.com/ldicarlo/legifrss/server/pkg/utils"
)

var consumerKey string
var consumerSecret string
var callbackURL string
var oauthToken string
var config oauth1.Config

func init() {
	envs, err := godotenv.Read(".env")
	if err != nil {
		panic("missing env file")
	}
	consumerKey = envs["twitter_consumer_key"]
	consumerSecret = envs["twitter_comsumer_secret"]
	callbackURL = envs["twitter_callback_url"]

	if consumerKey == "" || consumerSecret == "" || callbackURL == "" {
		panic("Missing one of the env params")
	}

}

func GetAuthURL() string {

	config = oauth1.Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		CallbackURL:    callbackURL,
		Endpoint:       oauthTwitter.AuthorizeEndpoint,
	}

	requestToken, _, err := config.RequestToken()
	utils.ErrCheck(err)
	authorizationURL, err := config.AuthorizationURL(requestToken)
	utils.ErrCheck(err)
	return authorizationURL.String()
}

func RegisterToken(newOauthToken string, tokenVerifier string) {
	oauthToken = newOauthToken
	token := oauth1.NewToken(oauthToken, tokenVerifier)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	tweets, resp, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: 20,
	})
	utils.ErrCheck(err)
	fmt.Println(tweets)
	fmt.Println(resp)

}

// TODO 1 => authenticate with twitter redirect, then save token locally

// TODO 2 => periodically publish tweets.
