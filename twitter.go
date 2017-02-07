package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var twitterClient *twitter.Client
var rawTwitterClient *http.Client

func randomTweet(s *discordgo.Session, m *discordgo.MessageCreate, query string) {
	search, _, err := twitterClient.Search.Tweets(&twitter.SearchTweetParams{Query: query, ResultType: "mixed"})
	if err != nil {
		log.Println(err)
		return
	}
	if len(search.Statuses) != 0 {
		s.ChannelMessageSend(m.ChannelID, "https://twitter.com/statuses/"+search.Statuses[rand.Intn(len(search.Statuses)-1)].IDStr)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Sadly there were no results for: "+query+" on twitter.")
	}
}
func trending(s *discordgo.Session, m *discordgo.MessageCreate) {
	type Trend []struct {
		Trends []struct {
			Name            string      `json:"name"`
			URL             string      `json:"url"`
			PromotedContent interface{} `json:"promoted_content"`
			Query           string      `json:"query"`
			TweetVolume     int         `json:"tweet_volume"`
		} `json:"trends"`
		AsOf      time.Time `json:"as_of"`
		CreatedAt time.Time `json:"created_at"`
		Locations []struct {
			Name  string `json:"name"`
			Woeid int    `json:"woeid"`
		} `json:"locations"`
	}
	resp, err := rawTwitterClient.Get(`https://api.twitter.com/1.1/trends/place.json?id=1`)
	if err != nil {
		log.Println("Failed to get Twitter trending topics:", err)
		return
	}
	defer resp.Body.Close()

	var record Trend
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
		return
	}

	var out []string
	for _, trends := range record {
		for _, trend := range trends.Trends {
			if trend.TweetVolume != 0 {
				out = append(out, fmt.Sprint("trend: ", trend.Name, " volume: ", trend.TweetVolume))
			}
		}
	}
	s.ChannelMessageSend(m.ChannelID, strings.Join(out, "\n"))
}
