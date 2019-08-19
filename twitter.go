package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
)

type twitterWrapper struct {
	*twitter.Client
	Raw *http.Client
}

const hanyuuDisplayName string = "Hanyuu_status"

var twitterClient twitterWrapper

func latestTweet() string {
	user, resp, err := twitterClient.Users.Show(&twitter.UserShowParams{
		ScreenName: hanyuuDisplayName,
	})
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalln(resp.Status)
	}
	return user.Status.Text
}

func getTweetID(link string) (int64, error) {
	u, err := url.Parse(link)
	if err != nil {
		return 0, err
	}

	u.RawQuery = "" //Drop parameters

	suffix := u.RequestURI()
	tokens := strings.Split(suffix, "/")
	for i, v := range tokens {
		if v == "status" && i < len(tokens)-1 {
			return strconv.ParseInt(tokens[i+1], 10, 64)
		}
	}
	return 0, errors.New("Invalid tweet URL")
}

// Attempts to grab images not embedded by discord and post them
func embedImages(s *discordgo.Session, m *discordgo.MessageCreate, link string) {
	ID, err := getTweetID(link)
	if err != nil {
		log.Println(err)
		channelMessageSendDeleteAble(s, m, "Sorry, this tweet URL seems malformed to me :(")
		return
	}
	// Pass the tweet ID to show, get tweet
	// https://godoc.org/github.com/dghubble/go-twitter/twitter#StatusService.Show
	statusShowParams := twitter.StatusShowParams{}
	statusShowParams.TweetMode = "extended"
	tweet, resp, err := twitterClient.Statuses.Show(ID, &statusShowParams)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(resp.Status)
		return
	}
	// Check if it has more than one image
	imageCount := 0
	msgResp := ""
	if tweet.ExtendedEntities == nil || tweet.ExtendedEntities.Media == nil {
		log.Println("Failed to read image data from tweet")
		return
	}
	for index, elem := range tweet.ExtendedEntities.Media {
		if elem.Type == "photo" {
			imageCount = imageCount + 1
		}
		if index != 0 { // If it's the first image, don't repost it
			msgResp += elem.MediaURL + "\n"
		}
	}
	// If so, post them!
	if imageCount > 1 {
		channelMessageSendDeleteAble(s, m, msgResp)
		return
	}
}

func randomTweet(s *discordgo.Session, m *discordgo.MessageCreate, query string) {
	search, _, err := twitterClient.Search.Tweets(&twitter.SearchTweetParams{Query: query})
	if err != nil {
		log.Println(err)
		return
	}
	if len(search.Statuses) == 0 {
		channelMessageSendDeleteAble(s, m, "Sadly there were no results for: "+query+" on twitter.")
		return
	}
	t := search.Statuses[rand.Intn(len(search.Statuses)-1)]
	if t.RetweetedStatus != nil {
		t = *t.RetweetedStatus
	}

	if _, err := channelMessageSendEmbedDeleteAble(s, m, tweetToEmbed(t)); err != nil {
		log.Println(err)
	}
}

func tweetToEmbed(t twitter.Tweet) *discordgo.MessageEmbed {
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
		log.Println("Failed to parse twitter time:", err.Error())
		return nil
	}
	embed := &discordgo.MessageEmbed{URL: "https://twitter.com/statuses/" + t.IDStr,
		Title: t.User.Name, Type: "rich", Timestamp: tim.Format(time.RFC3339Nano), Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Reweets: %d\tLikes: %d", t.RetweetCount, t.FavoriteCount)},
		Image: img, Thumbnail: thumb, Description: t.Text}
	embed.Fields = append(embed.Fields)

	return embed
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
	channelMessageSendDeleteAble(s, m, strings.Join(out, "\n"))
}
