package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/robfig/cron.v2"

	log "github.com/sirupsen/logrus"
)

var (
	client   *twitter.Client
	db       *sql.DB
	TweetOld []int64
)

type configStruct struct {
	ConsumerKey    string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
	AccessToken    string `json:"accessToken"`
	AccessSecret   string `json:"accessSecret"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	var TokenConfig configStruct
	user := os.Getenv("SQLUSER")
	pass := os.Getenv("SQLPASS")
	host := os.Getenv("DBHOST")

	if user == "" || pass == "" || host == "" {
		log.Error("user,pass or host not found")
		os.Exit(1)
	}

	file, err := ioutil.ReadFile("./token.json")

	if err != nil {
		log.Error(err)
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &TokenConfig)

	if err != nil {
		log.Error(err)
	}

	config := oauth1.NewConfig(TokenConfig.ConsumerKey, TokenConfig.ConsumerSecret)
	token := oauth1.NewToken(TokenConfig.AccessToken, TokenConfig.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client = twitter.NewClient(httpClient)
	db = Conn()
}

func main() {
	Group := GetGroup()
	for i := 0; i < len(Group); i++ {
		log.Info("Get old Data ", Group[i].VtuberGroup)
		tmp := GetTweetID(30, Group[i].ID)
		for j := 0; j < len(tmp); j++ {
			TweetOld = append(TweetOld, tmp[j])
		}
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	c := cron.New()
	c.AddFunc("@every 0h1m40s", CheckTweet)
	c.Start()
	wg.Wait()
}

func CheckTweet() {
	var TweetNew []int64
	Group := GetGroup()
	for i := 0; i < len(Group); i++ {
		log.Info("Check new Data ", Group[i].VtuberGroup)
		tmp := GetTweetID(30, Group[i].ID)
		for j := 0; j < len(tmp); j++ {
			TweetNew = append(TweetNew, tmp[j])
		}
	}
	for k := 0; k < len(TweetNew); k++ {
		if TweetNew[k] != TweetOld[k] {
			log.Info("New Tweet ", TweetNew[k])
			err := Like(TweetNew[k])
			if err != nil {
				log.Error(err)
				log.Info("Skip...")
			} else {
				log.Info("Done ", TweetNew[k])
				TweetOld = TweetNew
			}
		} else {
			log.Info("Still same")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func Like(twID int64) error {
	params := &twitter.FavoriteCreateParams{ID: twID}
	_, _, err := client.Favorites.Create(params)
	if err != nil {
		return err
	} else {
		err := Retweet(twID)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func Retweet(twID int64) error {
	_, _, err := client.Statuses.Retweet(twID, nil)
	if err != nil {
		return err
	}
	return nil
}
