package main

import (
	"github.com/bwmarrin/discordgo"
)

const appVersion = `12-05-2024
"Groq andy"`

func version(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "Current version: " + appVersion
	s.ChannelMessageSend(m.ChannelID, msg)
}
