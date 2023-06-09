package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Yuudachi's flavour replies.
func personality(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.Bot || m.Author.Username == "Liru" {
		return
	}

	// If the message is "poi" reply with "Poi!"
	pois := strings.Fields(strings.ToLower(m.Content))
	for _, word := range pois {
		if word == "poi" {
			ch, _ := s.Channel(m.ChannelID)
			gu, _ := s.Guild(ch.GuildID)
			for _, emoji := range gu.Emojis {
				if strings.Contains(strings.ToLower(emoji.Name), "poi") {
					if err := s.MessageReactionAdd(m.ChannelID, m.ID, emoji.APIName()); err != nil {
						log.Println(err)
					}
				}
			}
			return
		}
	}

	for _, user := range m.Mentions {
		if user.ID == botID {
			s.ChannelMessageSend(m.ChannelID, "Thank you for the kind message, <@"+m.Author.ID+">")
		}
	}
}
