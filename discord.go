package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const prefix string = "!"

func command(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == botID {
		return
	}
	//We don't like other bots either, Liru is an alias of a bot which didn't have an official bot account tag
	if m.Author.Bot || m.Author.Username == "Liru" {
		return
	}
	if len(m.Content) == 0 {
		return
	}

	if strings.HasPrefix(m.Content, prefix) && len(m.Content) > len(prefix) {
		m.Content = m.Content[len(prefix):]
		tokens := strings.Split(m.Content, " ")
		if tokens == nil {
			return
		}
		switch strings.ToLower(tokens[0]) {
		case "youtube":
			if len(tokens) > 1 {
				youtubeSearch(s, m, strings.Join(tokens[1:], " "))
			}
		case "next":
			nextVideo(s, m)
		case "limit":
			if len(tokens) > 1 {
				count, err := strconv.Atoi(tokens[1])
				if err != nil || count < 1 || count > 25 {
					s.ChannelMessageSend(m.ChannelID, "Failed to update limit")
					break
				}
				MaxYoutubeResults = count
			} else {
				s.ChannelMessageSend(m.ChannelID, "Current limit:"+fmt.Sprint(MaxYoutubeResults))
			}
		case "moon", "moonphase", "phase":
			moonPhase(s, m)
		case "version", "v":
			version(s, m)
		case "fortune", "f":
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
				var board = strings.Trim(tokens[1], `\/`)
				if BannedBoard(board) {
					s.ChannelMessageSend(m.ChannelID, "Sorry but this is a banned board")
					return
				}
				fourchan(s, m, tokens[1])
			} else {
				s.ChannelMessageSend(m.ChannelID, "Provide a board please!")
			}
		case "r", "radio", `r/a/dio`:
			if len(tokens) > 1 {
				switch tokens[1] {
				case "now", "current", "dj", "np":
					radioCurrent(s, m)
				case "q", "queue", "next":
					radioQueue(s, m)
				case "search", "s":
					if len(tokens) > 2 {
						radioSearch(s, m, strings.Join(tokens[2:], " "))
					}
				default:
					radioHelp(s, m)
				}
			} else {
				radioCurrent(s, m)
			}
		case "np", "song", "dj":
			radioCurrent(s, m)
		case "queue":
			radioQueue(s, m)
		case "b", "birb", "bird", "birds":
			birds(s, m)
		case "clap", "👏", "c":
			clap(s, m, tokens[1:])
		default:
		}
	}
}
