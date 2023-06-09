package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vharitonsky/iniflags"
)

var (
	botID string
)

var (
	discordBotToken = flag.String("token", "", "Discord Bot Token")
	youtubeAPIKey   = flag.String("youtube", "", "Youtube API Key with search permissions")
	printVersion    = flag.Bool("v", false, "Display current version")
)

func main() {
	iniflags.Parse()

	if *printVersion {
		fmt.Println(appVersion)
		os.Exit(0)
	}

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

	//here we add the functions
	dg.AddHandler(personality)
	dg.AddHandler(command)
	log.Println("Handlers added")

	if err := dg.Open(); err != nil {
		log.Fatalln("error opening connection,", err)
	}
	log.Println("Discord opened")

	timeOut := time.NewTicker(time.Second * 5)
	for {
		at := <-timeOut.C
		if err := dg.Open(); err != discordgo.ErrWSAlreadyOpen {
			log.Printf("%v occured at:%v\n", err, at)
		}
	}
}
