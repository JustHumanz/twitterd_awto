package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/robfig/cron.v2"

	log "github.com/sirupsen/logrus"
)

var (
	client *twitter.Client
	Data   TwitterD
	Auth   string
	Hana   string
	Holo   string
	Niji   string
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
	Auth = os.Getenv("AUTH")
	if Auth == "" {
		log.Error("Auth not found")
		os.Exit(1)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	c := cron.New()
	c.AddFunc("@every 0h1m20s", CheckTweetHanayori)
	c.AddFunc("@every 0h1m30s", CheckTweetNijisanji)
	c.AddFunc("@every 0h1m40s", CheckTweetHololive)
	c.Start()
	wg.Wait()
}

func CheckTweetHanayori() {
	log.Info("Start Curl Hanayori")
	body, err, _ := Curl("https://api.human-z.tech/vtbot/hanayori/twitter")
	if err != nil {
		log.Error(err)
	}
	err = json.Unmarshal(body, &Data)
	if err != nil {
		log.Error(err)
	}
	if Hana != Data[0].PermanentURL {
		for j := 0; j < 10; j++ {
			tmp, err := strconv.ParseInt(Data[j].TweetID, 10, 64)
			if err != nil {
				log.Error(err)
			}

			err = Like(tmp)
			if err != nil {
				log.Error(err)
				break
			} else {
				err = Retweet(tmp)
				if err != nil {
					log.Error(err)
				}
			}

		}

	} else {
		log.Info("Hanayori", " Still same")
	}
	Hana = Data[0].PermanentURL

}

func CheckTweetNijisanji() {
	log.Info("Start Curl Nijisanji")
	body, err, _ := Curl("https://api.human-z.tech/vtbot/nijisanji/twitter")
	if err != nil {
		log.Error(err)
	}
	err = json.Unmarshal(body, &Data)
	if err != nil {
		log.Error(err)
	}
	if Niji != Data[0].PermanentURL {
		for j := 0; j < 10; j++ {
			tmp, err := strconv.ParseInt(Data[j].TweetID, 10, 64)
			if err != nil {
				log.Error(err)
			}

			err = Like(tmp)
			if err != nil {
				log.Error(err)
				break
			} else {
				err = Retweet(tmp)
				if err != nil {
					log.Error(err)
				}
			}

		}
	} else {
		log.Info("Nijisanji", " Still same")
	}
	Niji = Data[0].PermanentURL
}

func CheckTweetHololive() {
	log.Info("Start Curl Hololive")
	body, err, _ := Curl("https://api.human-z.tech/vtbot/hololive/twitter")
	if err != nil {
		log.Error(err)
	}
	err = json.Unmarshal(body, &Data)
	if err != nil {
		log.Error(err)
	}
	if Holo != Data[0].PermanentURL {
		for j := 0; j < 10; j++ {
			tmp, err := strconv.ParseInt(Data[j].TweetID, 10, 64)
			if err != nil {
				log.Error(err)
			}

			err = Like(tmp)
			if err != nil {
				log.Error(err)
				break
			} else {
				err = Retweet(tmp)
				if err != nil {
					log.Error(err)
				}
			}

		}
	} else {
		log.Info("Hololive", " Still same")
	}
	Holo = Data[0].PermanentURL
}

func Like(twID int64) error {
	params := &twitter.FavoriteCreateParams{ID: twID}
	_, _, err := client.Favorites.Create(params)
	if err != nil {
		return err
	}
	log.Info("TweetID ", twID, " was liked")
	return nil
}

func Retweet(twID int64) error {
	_, _, err := client.Statuses.Retweet(twID, nil)
	if err != nil {
		return err
	}
	log.Info("TweetID ", twID, " was retweetd")
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
	req.Header.Set("Authorization", "Basic "+Auth)
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
