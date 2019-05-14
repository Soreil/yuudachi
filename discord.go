package main

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
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
		case "moon", "moonphase", "phase":
			moonPhase(s, m)
		case "delete", "delet":
			ChannelMessageDelete(s, m)
		case "version", "v":
			version(s, m)
		case "twitter":
			if len(tokens) > 1 {
				embedImages(s, m, tokens[1])
			}
		case "dog", "doggo", "goodboy":
			randomDogImage(s, m)
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
				switch strings.Trim(tokens[1], `\/`) {
				case "cm", "y", "gif", "e", "h", "hc", "b", "mlp", "lgbt", "soc", "s", "hm", "d", "t", "aco", "r", "pol", "trash":
					ChannelMessageSendDeleteAble(s, m, "I am a Christian bot, please don't make me blacklist you.\nFor now consider one of the following books instead for your reading pleasure.")
					bibleBooks(s, m)
					return
				}
				fourchan(s, m, tokens[1])
			} else {
				ChannelMessageSendDeleteAble(s, m, "Provide a board please!")
			}
		case "bible":
			if len(tokens) > 1 {
				bibleSearch(s, m, strings.Join(tokens[1:], " "))
			}
		case "breads", "bread", ":bread:", "🍞":
			if len(tokens) == 3 {
				n, err := strconv.ParseFloat(tokens[1], 64)
				if err != nil {
					ChannelMessageSendDeleteAble(s, m, "Failed to read conversion amount.")
					return
				}
				breads(s, m, n, strings.ToUpper(tokens[2]))
			}
			if len(tokens) == 4 && tokens[2] == "to" {
				n, err := strconv.ParseFloat(tokens[1], 64)
				if err != nil {
					ChannelMessageSendDeleteAble(s, m, "Failed to read conversion amount.")
					return
				}
				fiats(s, m, n, strings.ToUpper(tokens[3]))
			}
		case "r", "radio", `r/a/dio`, `r-a-dio`, `r-a-d.io`:
			if len(tokens) > 1 {
				switch tokens[1] {
				case "now", "current", "dj", "np":
					radioCurrent(s, m)
				case "q", "queue", "next":
					radioQueue(s, m)
				case "news":
					hanyuuUpdate(s, m)
				case "subscribe":
					radioSubscribe(s, m)
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
		case "queue", "next":
			radioQueue(s, m)
		case "b", "birb", "bird", "birds":
			birds(s, m)
		case "clap", "👏", "c":
			clap(s, m, tokens[1:])
		default:
			//ChannelMessageSendDeleteAble(s, m, "Unrecognized command: "+tokens[0])
			//usage(s, m)
		}
	}
}
