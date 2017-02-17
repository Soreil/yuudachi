package main

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"strings"
)

func command(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == botID {
		return
	}
	//We don't like other bots either
	if m.Author.Bot || m.Author.Username == "Liru" {
		return
	}
	if len(m.Content) == 0 {
		return
	}
	//We have a exclamation point
	if m.Content[0] == '!' && len(m.Content) > 2 {
		m.Content = m.Content[1:]
		tokens := strings.Split(m.Content, " ")
		if tokens == nil {
			return
		}
		switch strings.ToLower(tokens[0]) {
		case "twitter":
			if len(tokens) > 1 {
				switch tokens[1] {
				case "tweet", "search", "random":
					//Reuses the whole message
					randomTweet(s, m, strings.Join(tokens[2:], " "))
				case "trends", "trend", "trending":
					trending(s, m)
				}
			}
		case "version":
			version(s, m)
		case "fortune":
			if len(tokens) > 1 {
				//Only want one word since that's all the API can take.
				fortune(s, m, tokens[1])
			} else {
				//Can also be called without a word.
				fortune(s, m, "")
			}
		case "4chan", "4ch":
			if len(tokens) > 1 {
				//Only want one word since that's all the API can take.
				switch tokens[1] {
				case "cm", "y", "gif", "e", "h", "hc", "b", "mlp", "lgbt", "soc", "s", "hm", "d", "t", "aco", "r", "pol":
					s.ChannelMessageSend(m.ChannelID, "I am a Christian bot, please don't make me blacklist you.\nFor now consider one of the following books instead for your reading pleasure.")
					bibleBooks(s, m)
					return
				}
				fourchan(s, m, tokens[1])
			} else {
				s.ChannelMessageSend(m.ChannelID, "Provide a board please!")
			}
		case "bible":
			if len(tokens) > 1 {
				bibleSearch(s, m, strings.Join(tokens[1:], " "))
			}
		case "radio", `r/a/dio`, `r-a-dio`, `r-a-d.io`:
			if len(tokens) > 1 {
				//Only want one word since that's all the API can take.
				radio(s, m, tokens[1])
			} else {
				//Can also be called without a word.
				//fortune(s, m, "")
				radio(s, m, "")
			}
		case "8chan", "8ch":
			if len(tokens) > 1 {
				s.ChannelMessageSend(m.ChannelID, "I don't really like 8chan but maybe I'll let you look.")
				if roll := rand.Intn(6) + 1; roll != 6 {
					s.ChannelMessageSend(m.ChannelID, "You rolled a meagre: "+strconv.Itoa(roll)+"\nNo 8chan for you.")
					return
				}
				//Only want one word since that's all the API can take.
				eightchan(s, m, tokens[1])
			} else {
				s.ChannelMessageSend(m.ChannelID, "Provide a board please!")
			}
		case "birb", "bird", "birds":
			birds(s, m)
		default:
			s.ChannelMessageSend(m.ChannelID, "Unrecognized command: "+tokens[0])
			usage(s, m)
		}
	}
}
