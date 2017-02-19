package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func tatsumaki(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !(m.Author.Bot && m.Author.Username == "Tatsumaki") {
		return
	}
	if strings.Contains(strings.ToLower(m.Content), "leveled") {
		if err := s.MessageReactionAdd(m.ChannelID, m.ID, ":tada:"); err != nil {
			log.Println(err)
		}
	}
}
