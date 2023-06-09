package main

import (
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func birds(s *discordgo.Session, m *discordgo.MessageCreate) {
	//There isn't really a way to get a lot of all bird related emoji in unicode so for now it's just a hard coded list
	var birds = []string{"🦢", "🐥", "🐤", "🐣", "🐓", "🐔", "🐦", "🐧", "🕊️", "🦅", "🦆", "🦉", "🦚", "🦜"}

	//We want to use server specific fun old style emoji.
	ch, _ := s.Channel(m.ChannelID)
	gu, _ := s.Guild(ch.GuildID)
	for _, emoji := range gu.Emojis {
		if strings.ToLower(emoji.Name) == "birb" {
			if err := s.MessageReactionAdd(m.ChannelID, m.ID, emoji.APIName()); err != nil {
				log.Println(err)
			}
		}
	}

	for _, emoji := range gu.Emojis {
		if name := strings.ToLower(emoji.Name); strings.Contains(name, "birb") || strings.Contains(name, "bird") {
			birds = append(birds, "<:"+emoji.APIName()+">")
		}
	}

	s.ChannelMessageSend(m.ChannelID, birds[rand.Intn(len(birds))])
}
