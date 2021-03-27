package bot

import (
	"fmt"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	oauthTwitter "github.com/dghubble/oauth1/twitter"
	"github.com/joho/godotenv"
	"github.com/ldicarlo/legifrss/server/pkg/db"
	"github.com/ldicarlo/legifrss/server/pkg/utils"
)

var consumerKey string
var consumerSecret string
var callbackURL string
var twitterAccessTokenSecret string
var oauthToken string
var requestSecret string
var requestToken string
var config oauth1.Config

func init() {
	envs, err := godotenv.Read(".env")
	if err != nil {
		panic("missing env file")
	}
	consumerKey = envs["twitter_consumer_key"]
	consumerSecret = envs["twitter_comsumer_secret"]
	callbackURL = envs["twitter_callback_url"]
	twitterAccessTokenSecret = envs["twitter_access_token_secret"]

	if consumerKey == "" ||
		consumerSecret == "" ||
		callbackURL == "" ||
		twitterAccessTokenSecret == "" {
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

	rt, rs, err := config.RequestToken()
	utils.ErrCheck(err)
	requestToken = rt
	requestSecret = rs
	authorizationURL, err := config.AuthorizationURL(requestToken)
	utils.ErrCheck(err)
	return authorizationURL.String()
}

func RegisterToken(newOauthToken string, tokenVerifier string) {
	accessToken, accessSecret, err := config.AccessToken(requestToken, requestSecret, tokenVerifier)

	oauthToken = newOauthToken
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	fmt.Println(token)

	tweets, _, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: 20,
	})
	utils.ErrCheck(err)
	for _, tweet := range tweets {
		fmt.Println(tweet.Text)
	}

}

// TODO 1 => authenticate with twitter redirect, then save token locally

// TODO 2 => periodically publish tweets.
func GetElementsToPublish() {
	toPublish := db.ExtractContentToPublish()
	for _, elem := range toPublish {
		fmt.Println(elem.Title)
		elem.TwitterPublished = true
	}
	for _, elem := range toPublish {
		fmt.Println(elem.Title + ":" + strconv.FormatBool(elem.TwitterPublished))

	}
	db.Persist(toPublish)
}
