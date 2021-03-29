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
	toPublish, state := db.ExtractContentToPublish()
	for _, elem := range toPublish {
		if val, ok := state[elem.JORFID]; ok {
			result, err := PublishJORFAsTweets(elem, val)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Aborting")
				break
			}
			state[elem.JORFID] = result
		}
	}
	fmt.Println(state)
	db.PersistTwitterState(state)
}

func PublishJORFAsTweets(element models.JORFElement, twitterState models.TwitterJORF) (models.TwitterJORF, error) {
	if twitterState.StatusID == 0 {
		ID, err := publishJORFTweet(element)

		if !canContinue(err, element.JORFID) {
			return models.TwitterJORF{}, err
		}

		twitterState.StatusID = ID
	}
	for _, elem := range element.JORFContents {
		if twitterState.JORFContents[elem.ID] == 0 {
			statusID, err := publishLegifranceElementTweet(elem)
			if !canContinue(err, elem.ID) {
				return models.TwitterJORF{}, err
			}
			twitterState.JORFContents[elem.ID] = statusID
		}
	}

	return twitterState, nil
}

func canContinue(err error, jorfID string) bool {
	if err == nil {
		return true
	}
	if err.Error() == "twitter: 187 Status is a duplicate." {
		fmt.Println("Found duplicate, ignoring " + jorfID)
		return true
	} else {
		return false
	}
}

func publishJORFTweet(element models.JORFElement) (int64, error) {

	URI := ""
	if element.URI != "" {
		URI = " https://www.legifrance.gouv.fr" + element.URI
	}
	tweetStr := prepareTweetContent(element.JORFTitle) + URI
	// tweet, _, err := client.Statuses.Update(tweetStr, &twitter.StatusUpdateParams{})
	// return tweet.ID, err
	fmt.Println(element.JORFID + tweetStr)

	return 1, nil
}

func publishLegifranceElementTweet(element models.LegifranceElement) (int64, error) {

	tag := ""
	if element.Nature != "" {
		tag = " #" + element.Nature
	}
	tweetStr := prepareTweetContent(element.Description) + tag + " " + element.Link
	// tweet, _, err := client.Statuses.Update(tweetStr, &twitter.StatusUpdateParams{})
	// return tweet.ID, err
	fmt.Println("  " + element.ID + tweetStr)
	return 1, nil

}

func prepareTweetContent(str string) string {
	if len(str) < 200 {
		return str
	}
	return str[0:200] + "..."

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
