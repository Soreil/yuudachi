package main

import (
	"fmt"
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

	if m.Content[0] == '!' && len(m.Content) > 2 {
		m.Content = m.Content[1:]
		tokens := strings.Split(m.Content, " ")
		if tokens == nil {
			return
		}
		switch strings.ToLower(tokens[0]) {
		case "nightcore":
			if len(tokens) > 1 {
				tokens = append(tokens, "nightcore")
				youtubeSearch(s, m, strings.Join(tokens[1:], " "))
			}
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
					channelMessageSendDeleteAble(s, m, "Failed to update limit")
					break
				}
				MaxYoutubeResults = count
			} else {
				channelMessageSendDeleteAble(s, m, "Current limit:"+fmt.Sprint(MaxYoutubeResults))
			}
		case "moon", "moonphase", "phase":
			moonPhase(s, m)
		case "delete", "delet":
			ChannelMessageDeleteMostRecentOwnMessage(s, m)
		case "version", "v":
			version(s, m)
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
					channelMessageSendDeleteAble(s, m, "Sorry but this is a banned board")
					return
				}
				fourchan(s, m, tokens[1])
			} else {
				channelMessageSendDeleteAble(s, m, "Provide a board please!")
			}
		case "breads", "bread", ":bread:", "🍞":
			if len(tokens) == 3 {
				n, err := strconv.ParseFloat(tokens[1], 64)
				if err != nil {
					channelMessageSendDeleteAble(s, m, "Failed to read conversion amount.")
					return
				}
				breads(s, m, n, strings.ToUpper(tokens[2]))
			}
			if len(tokens) == 4 && tokens[2] == "to" {
				n, err := strconv.ParseFloat(tokens[1], 64)
				if err != nil {
					channelMessageSendDeleteAble(s, m, "Failed to read conversion amount.")
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
