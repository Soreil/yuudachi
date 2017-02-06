package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

const appVersion = "02-02-2017"

func version(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Println("Current version: " + appVersion)
	s.ChannelMessageSend(m.ChannelID, "Current version: "+appVersion)
}
