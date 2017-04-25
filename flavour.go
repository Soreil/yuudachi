package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

//Yuudachi's flavour replies.
func personality(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.Bot || m.Author.Username == "Liru" {
		if m.Author.Username == "Liru" && strings.Contains(strings.ToLower(m.Content), "awooo") {
			ChannelMessageSendDeleteAble(s, m, "No shouting Liru!")
		}
		return
	}

	// If the message is "poi" reply with "Poi!"
	pois := strings.Fields(strings.ToLower(m.Content))
	for _, word := range pois {
		if word == "poi" {
			_, err := ChannelMessageSendDeleteAble(s, m, "Poi!")
			if err != nil {
				log.Println(err)
			}
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

	if strings.Contains(strings.ToLower(m.Content), "kill "+strings.ToLower(botName)) {
		ChannelMessageSendDeleteAble(s, m, "EVASIVE MANOUVRES")
	}

	for _, user := range m.Mentions {
		if user.ID == botID {
			ChannelMessageSendDeleteAble(s, m, "Thank you for the kind message, <@"+m.Author.ID+">")
		}
	}
}
