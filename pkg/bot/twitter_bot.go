package bot

import (
	"fmt"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	oauthTwitter "github.com/dghubble/oauth1/twitter"
	"github.com/joho/godotenv"
	"github.com/ldicarlo/legifrss/server/pkg/db"
	"github.com/ldicarlo/legifrss/server/pkg/models"
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
var client *twitter.Client

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

	config = oauth1.Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		CallbackURL:    callbackURL,
		Endpoint:       oauthTwitter.AuthorizeEndpoint,
	}

	t, err := db.GetToken()
	if err == nil {
		fmt.Println("Found token")
		httpClient := config.Client(oauth1.NoContext, &t)
		client = twitter.NewClient(httpClient)

	}

}

func GetAuthURL() string {

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
	utils.ErrCheck(err)
	oauthToken = newOauthToken
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client = twitter.NewClient(httpClient)
	db.PersistToken(*token)
}

func ProcessElems() {
	toPublish := db.ExtractContentToPublish()
	var published []models.LegifranceElement
	for _, elem := range toPublish {
		fmt.Println(elem.Title)
		id, err := PublishElement(elem)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Aborting")
			break
		}
		elem.TwitterPublished = id
		published = append(published, elem)
	}
	db.Persist(published)
}

func PublishElement(element models.LegifranceElement) (int64, error) {
	var tweetStr string
	if len(element.Description) < 200 {
		tweetStr = element.Description
	} else {
		tweetStr = element.Description[0:200] + " ..."
	}
	tweet, _, err := client.Statuses.Update(tweetStr+" "+"#"+element.Nature+" "+element.Link, &twitter.StatusUpdateParams{})
	if err != nil {
		if err.Error() == "twitter: 187 Status is a duplicate." {
			fmt.Println("Found duplicate, ignoring " + element.ID)
			return 1, nil
		}
		return 0, err
	}
	return tweet.ID, nil
}

func GetAllTweets() []twitter.Tweet {
	tweets, _, _ := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: "Legifrance1",
		Count:      500,
	})
	return tweets
}

func RemoveAllTweets() []int64 {
	tweets := GetAllTweets()
	var result []int64
	for _, t := range tweets {
		fmt.Println("Destroying " + strconv.FormatInt(t.ID, 10))
		client.Statuses.Destroy(t.ID, &twitter.StatusDestroyParams{})
		result = append(result, t.ID)
	}
	return result
}
