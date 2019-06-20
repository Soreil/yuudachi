package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

//Clap simply does the emoji thing where we put a clapping emoji in between every word in a message.
func clap(s *discordgo.Session, m *discordgo.MessageCreate, c []string) {
	channelMessageSendDeleteAble(s, m, strings.Join(c, "👏"))
}
