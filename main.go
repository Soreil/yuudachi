package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/jaytaylor/html2text"
)

// Variables used for command line parameters
var (
	BotID         string
	Botname       string
	twitterClient *twitter.Client
	httpClient    *http.Client
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	discordToken := flags.String("token", "", "Discord Bot Token")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")
	flagutil.SetFlagsFromEnv(flags, "DISCORD")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" || *discordToken == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}
	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient = config.Client(oauth1.NoContext, token)

	// Twitter client
	twitterClient = twitter.NewClient(httpClient)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + *discordToken)
	if err != nil {
		log.Fatalln("error creating Discord session,", err)
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		log.Fatalln("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = u.ID
	Botname = u.Username
	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)
	// Register exclamation orchestrator as a callback for the sending events.
	dg.AddHandler(exclaim)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalln("error opening connection,", err)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
}

func exclaim(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}
	//We don't like other bots either
	if m.Author.Bot || m.Author.Username == "Liru" {
		return
	}
	if len(m.Content) == 0 {
		return
	}
	//We have a exclamation point
	if m.Content[0] == '!' {
		m.Content = m.Content[1:]
		tokens := strings.Split(m.Content, " ")
		if tokens == nil {
			return
		}
		switch tokens[0] {
		case "twitter":
			if len(tokens) > 1 {
				switch tokens[1] {
				case "tweet", "search", "random":
					//Reuses the whole message
					randomTweet(s, m, strings.Join(tokens[2:], " "))
				case "trends", "trend", "trending":
					trending(s, m)
				}
			}
		case "fortune":
			if len(tokens) > 1 {
				//Only want one word since that's all the API can take.
				fortune(s, m, tokens[1])
			} else {
				//Can also be called without a word.
				fortune(s, m, "")
			}
		case "4chan", "4ch":
			if len(tokens) > 1 {
				//Only want one word since that's all the API can take.
				fourchan(s, m, tokens[1])
			} else {
				s.ChannelMessageSend(m.ChannelID, "Provide a board please!")
			}
		case "radio", `r/a/dio`, `r-a-dio`, `r-a-d.io`:
			if len(tokens) > 1 {
				//Only want one word since that's all the API can take.
				radio(s, m, tokens[1])
			} else {
				//Can also be called without a word.
				//fortune(s, m, "")
				radio(s, m, "")
			}
		case "8chan", "8ch":
			if len(tokens) > 1 {
				//Only want one word since that's all the API can take.
				eightchan(s, m, tokens[1])
			} else {
				s.ChannelMessageSend(m.ChannelID, "Provide a board please!")
			}
		}
	}
}
func fourchan(s *discordgo.Session, m *discordgo.MessageCreate, board string) {
	type Catalog []struct {
		Page    int `json:"page"`
		Threads []struct {
			No            int    `json:"no"`
			Sticky        int    `json:"sticky,omitempty"`
			Closed        int    `json:"closed,omitempty"`
			Now           string `json:"now"`
			Name          string `json:"name"`
			Sub           string `json:"sub,omitempty"`
			Com           string `json:"com"`
			Filename      string `json:"filename"`
			Ext           string `json:"ext"`
			W             int    `json:"w"`
			H             int    `json:"h"`
			TnW           int    `json:"tn_w"`
			TnH           int    `json:"tn_h"`
			Tim           int64  `json:"tim"`
			Time          int    `json:"time"`
			Md5           string `json:"md5"`
			Fsize         int    `json:"fsize"`
			Resto         int    `json:"resto"`
			ID            string `json:"id"`
			Country       string `json:"country"`
			SemanticURL   string `json:"semantic_url"`
			CountryName   string `json:"country_name"`
			Replies       int    `json:"replies"`
			Images        int    `json:"images"`
			LastModified  int    `json:"last_modified"`
			Bumplimit     int    `json:"bumplimit,omitempty"`
			Imagelimit    int    `json:"imagelimit,omitempty"`
			OmittedPosts  int    `json:"omitted_posts,omitempty"`
			OmittedImages int    `json:"omitted_images,omitempty"`
			LastReplies   []struct {
				No          int    `json:"no"`
				Now         string `json:"now"`
				Name        string `json:"name"`
				Com         string `json:"com"`
				Time        int    `json:"time"`
				Resto       int    `json:"resto"`
				ID          string `json:"id"`
				Country     string `json:"country"`
				CountryName string `json:"country_name"`
				Filename    string `json:"filename,omitempty"`
				Ext         string `json:"ext,omitempty"`
				W           int    `json:"w,omitempty"`
				H           int    `json:"h,omitempty"`
				TnW         int    `json:"tn_w,omitempty"`
				TnH         int    `json:"tn_h,omitempty"`
				Tim         int64  `json:"tim,omitempty"`
				Md5         string `json:"md5,omitempty"`
				Fsize       int    `json:"fsize,omitempty"`
			} `json:"last_replies,omitempty"`
			MImg int `json:"m_img,omitempty"`
		} `json:"threads"`
	}

	url := "https://a.4cdn.org/"
	if board == "help" {
		s.ChannelMessageSend(m.ChannelID, "4chan usage:\n!4chan BOARD\nWhere category is one of the standard boards.")
		return
	}
	url += board + "/catalog.json"
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	if err == nil && resp.StatusCode != http.StatusOK {
		log.Println("Error: " + http.StatusText(resp.StatusCode))
		return
	}
	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record Catalog

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	page := record[rand.Intn(len(record)-1)]
	thread := page.Threads[rand.Intn(len(page.Threads)-1)]
	thread.Sub, err = html2text.FromString(strings.Replace(thread.Sub, "<wbr>", "", -1))
	if err != nil {
		log.Println(err)
		return
	}
	thread.Com, err = html2text.FromString(strings.Replace(thread.Com, "<wbr>", "", -1))
	if err != nil {
		log.Println(err)
		return
	}
	if strings.Contains(thread.Sub, "general") {
		log.Println("Error: We don't like generals.")
		fourchan(s, m, board)
		return
	}
	if thread.Sub != "" {
		thread.Sub = "*" + thread.Sub + "*" + "\n"
	}
	thread.Com = unembedURL(thread.Com)

	img := fmt.Sprintf("https://i.4cdn.org/%s/%d%s", board, thread.Tim, thread.Ext)
	link := fmt.Sprintf("<https://i.4cdn.org/%s/thread/%d>", board, thread.No)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s%s\n%s\n\n%s", thread.Sub, thread.Com, img, link))
}
func eightchan(s *discordgo.Session, m *discordgo.MessageCreate, board string) {
	type Catalog []struct {
		Threads []struct {
			No            int    `json:"no"`
			Sub           string `json:"sub,omitempty"`
			Com           string `json:"com"`
			Name          string `json:"name"`
			Trip          string `json:"trip,omitempty"`
			Time          int    `json:"time"`
			OmittedPosts  int    `json:"omitted_posts"`
			OmittedImages int    `json:"omitted_images"`
			Replies       int    `json:"replies"`
			Images        int    `json:"images"`
			Sticky        int    `json:"sticky"`
			Locked        int    `json:"locked"`
			Cyclical      string `json:"cyclical"`
			LastModified  int    `json:"last_modified"`
			TnH           int    `json:"tn_h"`
			TnW           int    `json:"tn_w"`
			H             int    `json:"h"`
			W             int    `json:"w"`
			Fsize         int    `json:"fsize"`
			Filename      string `json:"filename"`
			Ext           string `json:"ext"`
			Tim           string `json:"tim"`
			Resto         int    `json:"resto"`
			Country       string `json:"country,omitempty"`
			CountryName   string `json:"country_name,omitempty"`
			Md5           string `json:"md5,omitempty"`
			Email         string `json:"email,omitempty"`
		} `json:"threads"`
		Page int `json:"page"`
	}

	url := "https://8ch.net/"
	if board == "help" {
		s.ChannelMessageSend(m.ChannelID, "8chan usage:\n!8chan BOARD\nWhere category is one of the standard boards.")
		return
	}
	url += board + "/catalog.json"
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	if err == nil && resp.StatusCode != http.StatusOK {
		log.Println("Error: " + http.StatusText(resp.StatusCode))
		return
	}
	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record Catalog

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
		s, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(s))
	}

	page := record[rand.Intn(len(record))]
	thread := page.Threads[rand.Intn(len(page.Threads))]
	thread.Sub, err = html2text.FromString(strings.Replace(thread.Sub, "<wbr>", "", -1))
	if err != nil {
		log.Println(err)
		return
	}
	thread.Com, err = html2text.FromString(strings.Replace(thread.Com, "<wbr>", "", -1))
	if err != nil {
		log.Println(err)
		return
	}
	if strings.Contains(thread.Sub, "general") {
		log.Println("Error: We don't like generals.")
		eightchan(s, m, board)
		return
	}
	if thread.Sub != "" {
		thread.Sub = "*" + thread.Sub + "*" + "\n"
	}
	thread.Com = unembedURL(thread.Com)

	img := fmt.Sprintf("https://media.8ch.net/file_store/%s%s", thread.Tim, thread.Ext)
	req, err = http.NewRequest("GET", img, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client = &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	exist, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	if err == nil && exist.StatusCode != http.StatusOK {
		log.Println("Error: " + http.StatusText(exist.StatusCode))
		img = fmt.Sprintf("https://media.8ch.net/%s/src/%s%s", board, thread.Tim, thread.Ext)
	}
	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer exist.Body.Close()

	link := fmt.Sprintf("<https://8ch.net/%s/res/%d.html>", board, thread.No)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s%s\n%s\n\n%s", thread.Sub, thread.Com, img, link))
}

// Adds <> around the links to prevent embedding
func unembedURL(s string) string {
	start := strings.Index(s, "http")
	if start == -1 {
		return s
	}
	end := strings.IndexFunc(s[start:], unicode.IsSpace) + start
	if end > start {
		return s[:start] + "<" + s[start:end] + ">" + unembedURL(s[end:])
	}
	return s
}

func fortune(s *discordgo.Session, m *discordgo.MessageCreate, category string) {
	type Fortune struct {
		Fortune string `json:"fortune"`
	}
	url := "http://www.yerkee.com/api/fortune"
	if category == "help" {
		s.ChannelMessageSend(m.ChannelID, "Fortune usage:\n!fortune CATEGORY\nWhere category is one of: computers, cookie, definitions, miscellaneous, people, platitudes, politics, science, wisdom")
		return
	}

	if category == "computers" || category == "cookie" || category == "definitions" || category == "miscellaneous" || category == "people" || category == "platitudes" || category == "politics" || category == "science" || category == "wisdom" {
		url += "/" + category
	} else if category != "" {
		s.ChannelMessageSend(m.ChannelID, "Unknown category, type \"!fortune help\" for a list of categories allowed.")
		return
	}
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record Fortune

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	if record.Fortune != "" {
		s.ChannelMessageSend(m.ChannelID, record.Fortune)
	}
}

func randomTweet(s *discordgo.Session, m *discordgo.MessageCreate, query string) {
	search, _, err := twitterClient.Search.Tweets(&twitter.SearchTweetParams{Query: query, ResultType: "mixed"})
	if err != nil {
		log.Println(err)
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
	resp, err := httpClient.Get(`https://api.twitter.com/1.1/trends/place.json?id=1`)
	if err != nil {
		log.Fatalln("Failed to get Twitter trending topics:", err)
	}
	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record Trend
	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
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

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.Bot || m.Author.Username == "Liru" {
		return
	}

	// If the message is "poi" reply with "Poi!"
	if strings.Contains(strings.ToLower(m.Content), "poi") {
		_, err := s.ChannelMessageSendTTS(m.ChannelID, "Poi!")
		if err != nil {
			log.Println(err)
		}
		ch, _ := s.Channel(m.ChannelID)
		gu, _ := s.Guild(ch.GuildID)
		for _, emoji := range gu.Emojis {
			if strings.Contains(strings.ToLower(emoji.Name), "poi") {
				if err := s.MessageReactionAdd(m.ChannelID, m.ID, emoji.APIName()); err != nil {
					log.Println(err)
				}
			}
		}
	}

	//if message contains one of our emoji be thankful
	// if s.MessageReactionAdd
	// if strings.Contains(strings.ToLower(m.Content), "poi") {
	// 	_, err := s.ChannelMessageSend(m.ChannelID, "Poi!")
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }

	//SEND BOAT
	//	if strings.Contains(strings.ToLower(m.Content), "teleports away") {
	//		_, err := s.ChannelMessageSend(m.ChannelID, `https://upload.wikimedia.org/wikipedia/commons/8/8e/Yudachi_II.jpg`)
	//		if err != nil {
	//			log.Println(err)
	//		}
	//		return
	//	}

	if strings.Contains(strings.ToLower(m.Content), "kill "+strings.ToLower(Botname)) {
		s.ChannelMessageSend(m.ChannelID, "EVASIVE MANOUVRES")
	}
	// if m.Tts && m.Author.Bot || m.Author.Username == "Liru" {
	// 	//s.ChannelMessageSend(m.ChannelID, "No shouting Liru!")
	// 	_, err := s.ChannelMessageSendTTS(m.ChannelID, "Apoiiii")
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	for _, user := range m.Mentions {
		if user.ID == BotID {
			s.ChannelMessageSend(m.ChannelID, "Thank you for the kind message, <@"+m.Author.ID+">")
		}
	}
}

func radio(s *discordgo.Session, m *discordgo.MessageCreate, function string) {
	type Current struct {
		Main struct {
			Np           string `json:"np"`
			Listeners    int    `json:"listeners"`
			Bitrate      int    `json:"bitrate"`
			Isafkstream  bool   `json:"isafkstream"`
			Isstreamdesk bool   `json:"isstreamdesk"`
			Current      int    `json:"current"`
			StartTime    int    `json:"start_time"`
			EndTime      int    `json:"end_time"`
			Lastset      string `json:"lastset"`
			Trackid      int    `json:"trackid"`
			Thread       string `json:"thread"`
			Requesting   bool   `json:"requesting"`
			Djname       string `json:"djname"`
			Dj           struct {
				ID       int    `json:"id"`
				Djname   string `json:"djname"`
				Djtext   string `json:"djtext"`
				Djimage  string `json:"djimage"`
				Djcolor  string `json:"djcolor"`
				Visible  bool   `json:"visible"`
				Priority int    `json:"priority"`
				CSS      string `json:"css"`
				ThemeID  int    `json:"theme_id"`
				Role     string `json:"role"`
			} `json:"dj"`
			Queue []struct {
				Meta      string `json:"meta"`
				Time      string `json:"time"`
				Type      int    `json:"type"`
				Timestamp int    `json:"timestamp"`
			} `json:"queue"`
			Lp []struct {
				Meta      string `json:"meta"`
				Time      string `json:"time"`
				Type      int    `json:"type"`
				Timestamp int    `json:"timestamp"`
			} `json:"lp"`
		} `json:"main"`
	}
	// Ignore all messages created by the bot itself
	if m.Author.Bot || m.Author.Username == "Liru" {
		return
	}
	if function == "" {
		function = "dj"
	}
	url := `https://r-a-d.io/api`
	if function == "help" {
		s.ChannelMessageSend(m.ChannelID, "Usage: !radio dj")
		return
	}
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	if err == nil && resp.StatusCode != http.StatusOK {
		log.Println("Error: " + http.StatusText(resp.StatusCode))
		return
	}
	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record Current

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	var parts []string
	parts = append(parts, "Current DJ: "+record.Main.Dj.Djname)
	parts = append(parts, "Current song: "+record.Main.Lp[0].Meta)
	if record.Main.Thread != "" {
		parts = append(parts, "There is a thread up: "+record.Main.Thread)

	}
	parts = append(parts, "https://r-a-d.io/api/dj-image/"+record.Main.Dj.Djimage)

	s.ChannelMessageSend(m.ChannelID, strings.Join(parts, "\n"))
	return
}
