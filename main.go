package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"strconv"
)

// Variables used for command line parameters
var (
	botID   string
	botName string
)

const userAgent = "Yuudachi/0.1"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	twitterConsumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	twitterConsumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	twitterAccessToken := flags.String("access-token", "", "Twitter Access Token")
	twitterAccessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	discordBotToken := flags.String("token", "", "Discord Bot Token")
	bibleAccessToken := flags.String("bible", "", "Bible search token")
	printVersion := flags.Bool("v", false, "Display current version")
	flags.Parse(os.Args[1:])

	//TODO(Sjon): Give this a proper place.
	if *printVersion {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	if *twitterConsumerKey == "" || *twitterConsumerSecret == "" || *twitterAccessToken == "" || *twitterAccessSecret == "" || *discordBotToken == "" || *bibleAccessToken == "" {
		log.Println(*twitterConsumerKey, *twitterConsumerSecret, *twitterAccessToken, *twitterAccessSecret, *discordBotToken, *bibleAccessToken)
		log.Fatal("Consumer key/secret and Access token/secret required")
	}
	log.Println("Keys gotten")
	bibleToken = *bibleAccessToken
	config := oauth1.NewConfig(*twitterConsumerKey, *twitterConsumerSecret)
	token := oauth1.NewToken(*twitterAccessToken, *twitterAccessSecret)
	log.Println("Twitter tokens done")
	// OAuth1 http.Client will automatically authorize Requests
	rawTwitterClient = config.Client(oauth1.NoContext, token)

	// Twitter client
	twitterClient = twitter.NewClient(rawTwitterClient)
	log.Println("Twitter set up")
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + *discordBotToken)
	if err != nil {
		log.Fatalln("error creating Discord session,", err)
	}
	log.Println("Discord session created")

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		log.Fatalln("error obtaining account details,", err)
	}
	log.Println("Got bot details")

	botID = u.ID
	botName = u.Username
	dg.AddHandler(personality)
	dg.AddHandler(command)
	log.Println("Handlers added")
	err = dg.Open()
	if err != nil {
		log.Fatalln("error opening connection,", err)
	}
	log.Println("Discord opened")
	fmt.Println("Succesfully initialized")
	<-make(chan struct{})
}

func command(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == botID {
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
		switch strings.ToLower(tokens[0]) {
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
		case "version":
			version(s, m)
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
				switch tokens[1] {
				case "cm", "y", "gif", "e", "h", "hc", "b", "mlp", "lgbt", "soc", "s", "hm", "d", "t", "aco", "r", "pol":
					s.ChannelMessageSend(m.ChannelID, "I am a Christian bot, please don't make me blacklist you.\nFor now consider one of the following books instead for your reading pleasure.")
					bibleBooks(s, m)
					return
				}
				fourchan(s, m, tokens[1])
			} else {
				s.ChannelMessageSend(m.ChannelID, "Provide a board please!")
			}
		case "bible":
			if len(tokens) > 1 {
				bibleSearch(s, m, strings.Join(tokens[1:], " "))
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
				s.ChannelMessageSend(m.ChannelID, "I don't really like 8chan but maybe I'll let you look.")
				if roll := rand.Intn(6) + 1; roll != 6 {
					s.ChannelMessageSend(m.ChannelID, "You rolled a meagre: "+strconv.Itoa(roll)+"\nNo 8chan for you.")
					return
				}
				//Only want one word since that's all the API can take.
				eightchan(s, m, tokens[1])
			} else {
				s.ChannelMessageSend(m.ChannelID, "Provide a board please!")
			}
		case "birb", "bird", "birds":
			birds(s, m)
		default:
			s.ChannelMessageSend(m.ChannelID, "Unrecognized command: "+tokens[0])
			usage(s, m)
		}
	}
}
