package main

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

var posts struct {
	m map[*discordgo.Guild]map[string][]*discordgo.Message
	sync.Mutex
}

func init() {
	posts.m = make(map[*discordgo.Guild]map[string][]*discordgo.Message)
}

func ChannelMessageDelete(s *discordgo.Session, m *discordgo.MessageCreate) {
	ch, _ := s.Channel(m.ChannelID)
	gu, _ := s.Guild(ch.GuildID)

	posts.Lock()
	defer posts.Unlock()

	err := s.ChannelMessageDelete(posts.m[gu][ch.ID][len(posts.m[gu])-1].ChannelID, posts.m[gu][ch.ID][len(posts.m[gu])-1].ID)
	posts.m[gu][ch.ID] = posts.m[gu][ch.ID][:len(posts.m[gu])-1]
	if err != nil {
		panic(err)
	}
}

func addChannelMessageDeleteAble(s *discordgo.Session, m *discordgo.MessageCreate, message *discordgo.Message) {
	ch, _ := s.Channel(m.ChannelID)
	gu, _ := s.Guild(ch.GuildID)
	if _, ok := posts.m[gu]; !ok {
		posts.m[gu] = make(map[string][]*discordgo.Message)
	}
	posts.m[gu][ch.ID] = append(posts.m[gu][ch.ID], message)

}

func ChannelMessageSendDeleteAble(s *discordgo.Session, m *discordgo.MessageCreate, content string) (*discordgo.Message, error) {
	posts.Lock()
	defer posts.Unlock()

	message, err := s.ChannelMessageSend(m.ChannelID, content)
	addChannelMessageDeleteAble(s, m, message)
	return message, err
}

func _(s *discordgo.Session, m *discordgo.MessageCreate, content string) (*discordgo.Message, error) {
	posts.Lock()
	defer posts.Unlock()
	message, err := s.ChannelMessageSendTTS(m.ChannelID, content)
	addChannelMessageDeleteAble(s, m, message)
	return message, err
}

func ChannelMessageSendEmbedDeleteAble(s *discordgo.Session, m *discordgo.MessageCreate, content *discordgo.MessageEmbed) (*discordgo.Message, error) {
	posts.Lock()
	defer posts.Unlock()
	message, err := s.ChannelMessageSendEmbed(m.ChannelID, content)
	addChannelMessageDeleteAble(s, m, message)
	return message, err
}
