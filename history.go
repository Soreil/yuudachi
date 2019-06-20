package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var posts struct {
	m map[*discordgo.Guild]map[string][]*discordgo.Message
	sync.Mutex
}

func init() {
	posts.m = make(map[*discordgo.Guild]map[string][]*discordgo.Message)
}

//mostRecentPost returns the message ID for the most recent post on the current server and channel
func mostRecentPost(s *discordgo.Session, m *discordgo.MessageCreate) *discordgo.Message {
	ch, _ := s.Channel(m.ChannelID)
	gu, _ := s.Guild(ch.GuildID)

	posts.Lock()
	defer posts.Unlock()

	if _, ok := posts.m[gu]; !ok {
		posts.m[gu] = make(map[string][]*discordgo.Message)
	}

	messages := posts.m[gu][ch.ID]

	if len(messages) == 0 {
		return nil
	}

	return messages[len(messages)-1]
}

//ChannelMessageDeleteMostRecent is a wrapper around ChannelMessageDelete
func ChannelMessageDeleteMostRecent(s *discordgo.Session, m *discordgo.MessageCreate) {
	recent := mostRecentPost(s, m)
	if recent == nil {
		channelMessageSendDeleteAble(s, m, "No messages to delete left in this channel :)")
		return
	}

	//We have removed the post corresponding to mID
	err := s.ChannelMessageDelete(recent.ChannelID, recent.ID)
	if err != nil {
		panic(err)
	}
	if err := removeChannelMessageFromQueue(s, m, recent); err != nil {
		panic(err)
	}

}

func removeChannelMessageFromQueue(s *discordgo.Session, m *discordgo.MessageCreate, message *discordgo.Message) error {

	ch, _ := s.Channel(message.ChannelID)
	gu, _ := s.Guild(ch.GuildID)

	posts.Lock()
	defer posts.Unlock()

	if _, ok := posts.m[gu]; !ok {
		posts.m[gu] = make(map[string][]*discordgo.Message)
	}

	messages := posts.m[gu][ch.ID]
	for i := range messages {
		if messages[i].ID == message.ID {
			posts.m[gu][ch.ID] = append(messages[:i], messages[i+1:]...)
			return nil
		}
	}
	return errors.New("Failed to remove message: " + fmt.Sprint(message))
}

func addChannelMessageDeleteAble(s *discordgo.Session, m *discordgo.MessageCreate, message *discordgo.Message) {
	ch, _ := s.Channel(m.ChannelID)
	gu, _ := s.Guild(ch.GuildID)

	if _, ok := posts.m[gu]; !ok {
		posts.m[gu] = make(map[string][]*discordgo.Message)
	}
	posts.m[gu][ch.ID] = append(posts.m[gu][ch.ID], message)

}

func channelMessageSendDeleteAble(s *discordgo.Session, m *discordgo.MessageCreate, content string) (*discordgo.Message, error) {
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

func channelMessageSendEmbedDeleteAble(s *discordgo.Session, m *discordgo.MessageCreate, content *discordgo.MessageEmbed) (*discordgo.Message, error) {
	posts.Lock()
	defer posts.Unlock()
	message, err := s.ChannelMessageSendEmbed(m.ChannelID, content)
	addChannelMessageDeleteAble(s, m, message)
	return message, err
}
