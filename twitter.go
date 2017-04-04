package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type twitterWrapper struct {
	*twitter.Client
	Raw *http.Client
}

var twitterClient twitterWrapper

func latestTweet() {
	p := twitter.StreamFilterParams{Follow: []string{`@Hanyuu_status`}}
	stream, err := twitterClient.Streams.Filter(&p)
	if err != nil {
		log.Println("Failed to connect to Twitter Streaming API", p)
		return
	}
	for msg := range stream.Messages {
		fmt.Println("Trying to figure out the type of a Twitter API response")
		fmt.Println(msg)
		fmt.Println(reflect.TypeOf(msg))
	}
	defer stream.Stop()
}

func randomTweet(s *discordgo.Session, m *discordgo.MessageCreate, query string) {
	search, _, err := twitterClient.Search.Tweets(&twitter.SearchTweetParams{Query: query})
	if err != nil {
		log.Println(err)
		return
	}
	if len(search.Statuses) != 0 {
		t := search.Statuses[rand.Intn(len(search.Statuses)-1)]
		if t.RetweetedStatus != nil {
			t = *t.RetweetedStatus
		}
		//log.Printf("%+v\n",t.Entities)
		//log.Printf("%+v\n",t.ExtendedEntities)

		img := &discordgo.MessageEmbedImage{}
		thumb := &discordgo.MessageEmbedThumbnail{}

		if len(t.Entities.Media) > 0 {
			img.URL = t.Entities.Media[0].MediaURLHttps
			img.Height = t.Entities.Media[0].Sizes.Medium.Height
			img.Width = t.Entities.Media[0].Sizes.Medium.Height

		}
		thumb.URL = t.User.ProfileImageURLHttps

		tim, err := time.Parse(time.RubyDate, t.CreatedAt)
		if err != nil {
			log.Printf("Failed to parse twitter time: %s", err.Error())
		}
		embed := &discordgo.MessageEmbed{URL: "https://twitter.com/statuses/" + t.IDStr,
			Title: t.User.Name, Type: "rich", Timestamp: tim.Format(time.RFC3339Nano), Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Reweets: %d\tLikes: %d", t.RetweetCount, t.FavoriteCount)},
			Image: img, Thumbnail: thumb, Description: t.Text}
		embed.Fields = append(embed.Fields)
		if _, err := s.ChannelMessageSendEmbed(m.ChannelID, embed); err != nil {
			log.Println(err, tim.Format(RFC3339Discord), tim.Format(time.RFC3339Nano), tim, t.CreatedAt)
		}
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
	resp, err := twitterClient.Raw.Get(`https://api.twitter.com/1.1/trends/place.json?id=1`)
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
