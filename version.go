package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

const appVersion = "07-02-2017"

func version(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Println("Current version: " + appVersion)
	s.ChannelMessageSend(m.ChannelID, "Current version: "+appVersion)
}

func usage(s *discordgo.Session, m *discordgo.MessageCreate) {
	usage := strings.Join([]string{"twitter", "version", "fortune", "8chan", "4chan", "bible", "radio", "bird"}, ", ")
	s.ChannelMessageSend(m.ChannelID, "The possible commands Yuudachi will like: "+usage+".")
}
