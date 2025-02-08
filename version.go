package main

import (
	"github.com/bwmarrin/discordgo"
)

const appVersion = `08-02-2025
"Chinese"`

func version(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "Current version: " + appVersion
	s.ChannelMessageSend(m.ChannelID, msg)
}
