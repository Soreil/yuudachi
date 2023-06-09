package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Clap simply does the emoji thing where we put a clapping emoji in between every word in a message.
func clap(s *discordgo.Session, m *discordgo.MessageCreate, c []string) {
	s.ChannelMessageSend(m.ChannelID, strings.Join(c, "👏"))
}
