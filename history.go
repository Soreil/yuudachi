package main

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type messageRelationship struct {
	cause *discordgo.Message
	reply *discordgo.Message
}

var posts struct {
	m map[*discordgo.Guild]map[string][]messageRelationship
	sync.Mutex
}

func init() {
	posts.m = make(map[*discordgo.Guild]map[string][]messageRelationship)
}

//mostRecentPost returns the message ID for the most recent post on the current server and channel
func mostRecentOwnPost(s *discordgo.Session, m *discordgo.MessageCreate) messageRelationship {
	ch, _ := s.Channel(m.ChannelID)
	gu, _ := s.Guild(ch.GuildID)

	posts.Lock()
	defer posts.Unlock()

	if _, ok := posts.m[gu]; !ok {
		posts.m[gu] = make(map[string][]messageRelationship)
	}

	messages := posts.m[gu][ch.ID]
	if messages == nil {
		log.Println("Warning: we got an empty message array, this should never happen")
		return messageRelationship{}
	}

	for i := range messages {
		msg := messages[len(messages)-i-1]
		if msg.cause.Author.ID == m.Author.ID {
			return msg
		}
	}
	return messageRelationship{}
}

//ChannelMessageDeleteMostRecentOwnMessage is a wrapper around ChannelMessageDelete
func ChannelMessageDeleteMostRecentOwnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	recent := mostRecentOwnPost(s, m)
	if recent.reply == nil {
		channelMessageSendDeleteAble(s, m, "No messages to delete left in this channel :)")
		return
	}

	//We have removed the post corresponding to mID
	err := s.ChannelMessageDelete(recent.reply.ChannelID, recent.reply.ID)
	if err != nil {
		panic(err)
	}
	if err := removeChannelMessageFromQueue(s, m, recent); err != nil {
		panic(err)
	}

}
func removeChannelMessageFromQueue(s *discordgo.Session, m *discordgo.MessageCreate, message messageRelationship) error {

	ch, _ := s.Channel(message.reply.ChannelID)
	gu, _ := s.Guild(ch.GuildID)

	posts.Lock()
	defer posts.Unlock()

	if _, ok := posts.m[gu]; !ok {
		posts.m[gu] = make(map[string][]messageRelationship)
	}

	messages := posts.m[gu][ch.ID]
	for i := range messages {
		if messages[i].reply.ID == message.reply.ID {
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
		posts.m[gu] = make(map[string][]messageRelationship)
	}
	posts.m[gu][ch.ID] = append(posts.m[gu][ch.ID], messageRelationship{cause: m.Message, reply: message})

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
