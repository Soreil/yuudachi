package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const appVersion = `14-5-2019
"we now have a basic twitter media reposter for images not embedded by Discord"`

func version(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "Current version: " + appVersion
	ChannelMessageSendDeleteAble(s, m, msg)
}

func usage(s *discordgo.Session, m *discordgo.MessageCreate) {
	usage := strings.Join([]string{"twitter", "version", "fortune", "4chan", "bible", "radio", "bird"}, ", ")
	ChannelMessageSendDeleteAble(s, m, "The possible commands Yuudachi will like: "+usage+".")
}
