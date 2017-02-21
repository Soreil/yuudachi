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
		if err := s.MessageReactionAdd(m.ChannelID, m.ID, "🎉"); err != nil {
			log.Println("Failed to to add 🎉", err)
		}
	}
	if strings.Contains(strings.ToLower(m.Content), "gelbooru") {
		if _, err := s.ChannelMessageSend(m.ChannelID, clear); err != nil {
			log.Println(err)
		}
	}
}

const clear = `





























































































































































`
