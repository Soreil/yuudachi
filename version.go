package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const appVersion = `07-01-2020
"Added back in bread support :)"`

func version(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "Current version: " + appVersion
	channelMessageSendDeleteAble(s, m, msg)
}

func usage(s *discordgo.Session, m *discordgo.MessageCreate) {
	usage := strings.Join([]string{"twitter", "version", "fortune", "4chan", "bible", "radio", "bird"}, ", ")
	channelMessageSendDeleteAble(s, m, "The possible commands Yuudachi will like: "+usage+".")
}
