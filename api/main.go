package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/robfig/cron.v2"

	log "github.com/sirupsen/logrus"
)

var (
	client    *twitter.Client
	TweetOld  []int64
	GroupName []string
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
}

func main() {
	log.Info("Start bot")
	GroupName = []string{"hanayori", "nijisanji", "hololive"}
	for i := 0; i < len(GroupName); i++ {
		var Data TwitterD

		log.Info("Start Curl " + GroupName[i])
		body, err, _ := Curl("https://api.justhumanz.me/BotAPI/" + GroupName[i] + "/twitter")
		if err != nil {
			log.Error(err)
		}
		err = json.Unmarshal(body, &Data)
		if err != nil {
			log.Error(err)
		}
		for j := 0; j < len(Data); j++ {
			tmp, err := strconv.ParseInt(Data[j].TweetID, 10, 64)
			if err != nil {
				log.Error(err)
			}
			TweetOld = append(TweetOld, tmp)
		}
	}
	CheckTweet()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	c := cron.New()
	c.AddFunc("@every 0h1m40s", CheckTweet)
	c.Start()
	wg.Wait()
}

func CheckTweet() {
	log.Info("Start Cron")
	var TweetNew []int64
	for i := 0; i < len(GroupName); i++ {
		var Data TwitterD
		body, err, _ := Curl("https://api.justhumanz.me/BotAPI/" + GroupName[i] + "/twitter")
		if err != nil {
			log.Error(err)
		}
		err = json.Unmarshal(body, &Data)
		if err != nil {
			log.Error(err)
		}
		for j := 0; j < len(Data); j++ {
			tmp, err := strconv.ParseInt(Data[j].TweetID, 10, 64)
			if err != nil {
				log.Error(err)
			}
			TweetNew = append(TweetNew, tmp)
		}

	}
	log.Info("Len : ", len(TweetNew))

	for i := 0; i < len(TweetNew); i++ {
		if TweetNew[i] != TweetOld[i] {
			log.WithFields(log.Fields{
				"TweetID": TweetNew[i],
			}).Info("New Post")
			err := Like(TweetNew[i])
			if err != nil {
				log.Error(err)
			} else {
				err = Retweet(TweetNew[i])
				if err != nil {
					log.Error(err)
				}
			}
		} else {
			log.WithFields(log.Fields{
				"TweetID": TweetNew[i],
			}).Info("Still same")
		}
	}
	TweetOld = TweetNew
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

func Curl(url string) ([]byte, error, int) {
	spaceClient := http.Client{
		Timeout: time.Second * 20, // Timeout after 20 seconds
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)
		return []byte{}, err, 0
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Nokia 3310) AppleWebKit/601.2 (KHTML, like Gecko)")
	req.Header.Set("Authorization", "Basic a2FubzprYW5vMjUyNQ==")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Error(getErr)
		return []byte{}, err, res.StatusCode
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Error(readErr)

	}
	return body, nil, res.StatusCode
}

type TwitterD []struct {
	VtuberName      string   `json:"VtuberName"`
	VtuberNameEN    string   `json:"VtuberName_EN"`
	VtuberNameJP    string   `json:"VtuberName_JP"`
	VtuberGroupName string   `json:"VtuberGroupName"`
	PermanentURL    string   `json:"PermanentURL"`
	Author          string   `json:"Author"`
	Likes           int      `json:"Likes"`
	Photos          []string `json:"Photos"`
	Videos          string   `json:"Videos"`
	Text            string   `json:"Text"`
	TweetID         string   `json:"TweetID"`
}
