package main

import (
	"github.com/bwmarrin/discordgo"
)

const appVersion = `18-03-2021
"no more rolls Sadge"`

func version(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "Current version: " + appVersion
	channelMessageSendDeleteAble(s, m, msg)
}
