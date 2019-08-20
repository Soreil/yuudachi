package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/vharitonsky/iniflags"
)

var (
	botID   string
	botName string
)

var (
	twitterConsumerKey    = flag.String("consumer-key", "", "Twitter Consumer Key")
	twitterConsumerSecret = flag.String("consumer-secret", "", "Twitter Consumer Secret")
	twitterAccessToken    = flag.String("access-token", "", "Twitter Access Token")
	twitterAccessSecret   = flag.String("access-secret", "", "Twitter Access Secret")
	discordBotToken       = flag.String("token", "", "Discord Bot Token")
	bibleAccessToken      = flag.String("bible", "", "Bible search token")
	fixerAPIToken         = flag.String("fixer", "", "Fixer currency token")
	youtubeAPIKey         = flag.String("youtube", "", "Youtube API Key with search permissions")
	printVersion          = flag.Bool("v", false, "Display current version")
)

const userAgent = "Yuudachi/0.2"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	iniflags.Parse()

	if *printVersion {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	if *twitterConsumerKey == "" || *twitterConsumerSecret == "" || *twitterAccessToken == "" || *twitterAccessSecret == "" || *discordBotToken == "" || *bibleAccessToken == "" || *youtubeAPIKey == "" {
		log.Println(*twitterConsumerKey, *twitterConsumerSecret, *twitterAccessToken, *twitterAccessSecret, *discordBotToken, *bibleAccessToken)
		log.Fatal("Consumer key/secret and Access token/secret required")
	}
	log.Println("Keys gotten")
	bibleToken = *bibleAccessToken

	config := oauth1.NewConfig(*twitterConsumerKey, *twitterConsumerSecret)
	token := oauth1.NewToken(*twitterAccessToken, *twitterAccessSecret)
	log.Println("Twitter tokens done")
	// OAuth1 http.Client will automatically authorize Requests
	twitterClient = twitterWrapper{nil, config.Client(oauth1.NoContext, token)}
	twitterClient.Client = twitter.NewClient(twitterClient.Raw)

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

	//here we add the functions
	dg.AddHandler(personality)
	dg.AddHandler(command)
	log.Println("Handlers added")

	if err := dg.Open(); err != nil {
		log.Fatalln("error opening connection,", err)
	}
	log.Println("Discord opened")

	fmt.Println("Succesfully initialized", botName)
	<-make(chan struct{})
}
