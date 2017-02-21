package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"strings"
)

func birds(s *discordgo.Session, m *discordgo.MessageCreate) {
	var birds = []string{":bird:", ":dove:", ":chicken:", ":baby_chick:", ":rooster:", ":penguin:", ":turkey:", ":eagle:", ":duck:", ":owl:"}

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
	s.ChannelMessageSend(m.ChannelID, birds[rand.Intn(len(birds))])
}
