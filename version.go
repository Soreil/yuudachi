package main

import (
	"github.com/bwmarrin/discordgo"
)

const appVersion = `11-05-2024
"Holy moly"`

func version(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "Current version: " + appVersion
	s.ChannelMessageSend(m.ChannelID, msg)
}
