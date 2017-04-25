package main

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

var posts struct {
	m map[*discordgo.Guild][]*discordgo.Message
	sync.Mutex
}

func init() {
	posts.m = make(map[*discordgo.Guild][]*discordgo.Message)
}

func ChannelMessageDelete(s *discordgo.Session, m *discordgo.MessageCreate) {
	ch, _ := s.Channel(m.ChannelID)
	gu, _ := s.Guild(ch.GuildID)

	posts.Lock()
	defer posts.Unlock()

	if len(posts.m[gu]) == 0 {
		return
	}
	err := s.ChannelMessageDelete(posts.m[gu][len(posts.m[gu])-1].ChannelID, posts.m[gu][len(posts.m[gu])-1].ID)
	posts.m[gu] = posts.m[gu][:len(posts.m[gu])-1]
	if err != nil {
		panic(err)
	}
}

func ChannelMessageSendDeleteAble(s *discordgo.Session, m *discordgo.MessageCreate, content string) (*discordgo.Message, error) {
	posts.Lock()
	defer posts.Unlock()

	message, err := s.ChannelMessageSend(m.ChannelID, content)
	ch, _ := s.Channel(m.ChannelID)
	gu, _ := s.Guild(ch.GuildID)
	posts.m[gu] = append(posts.m[gu], message)

	return message, err
}

func ChannelMessageSendTTSDeleteAble(s *discordgo.Session, m *discordgo.MessageCreate, content string) (*discordgo.Message, error) {
	posts.Lock()
	defer posts.Unlock()
	message, err := s.ChannelMessageSendTTS(m.ChannelID, content)
	ch, _ := s.Channel(m.ChannelID)
	gu, _ := s.Guild(ch.GuildID)
	posts.m[gu] = append(posts.m[gu], message)

	return message, err
}

func ChannelMessageSendEmbedDeleteAble(s *discordgo.Session, m *discordgo.MessageCreate, content *discordgo.MessageEmbed) (*discordgo.Message, error) {
	posts.Lock()
	defer posts.Unlock()
	message, err := s.ChannelMessageSendEmbed(m.ChannelID, content)
	ch, _ := s.Channel(m.ChannelID)
	gu, _ := s.Guild(ch.GuildID)
	posts.m[gu] = append(posts.m[gu], message)

	return message, err
}
