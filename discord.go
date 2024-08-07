package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const prefix string = "!"

var groqLUT map[string][]Message = map[string][]Message{}

func command(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == botID {
		return
	}
	//We don't like other bots either
	if m.Author.Bot {
		return
	}
	if len(m.Content) == 0 {
		return
	}

	if m.ReferencedMessage != nil {
		ctx, contains := groqLUT[m.ReferencedMessage.ID]
		if contains {
			msg, history, err := AskGroq(m.Content, ctx)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Oopsie woopsie we got an error groq sisters:"+err.Error())
				log.Println(err)
			}

			var messages = chunk(msg, 1000)
			for _, v := range messages {
				msg, err := s.ChannelMessageSend(m.ChannelID, v)
				if err != nil {
					panic(err)
				}
				groqLUT[msg.ID] = history
			}
			return
		}
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
		case "groq":

			msg, history, err := func() (string, []Message, error) {
				if m.ReferencedMessage != nil {
					var ctx = groqLUT[m.ReferencedMessage.ID]
					return AskGroq(strings.Join(tokens[1:], " "), ctx)
				} else {
					var x = AskGroqSystem("You are a discord bot which has to answer questions succinctly as to not flood the channel", nil)
					return AskGroq(strings.Join(tokens[1:], " "), x)
				}
			}()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Oopsie woopsie we got an error groq sisters:"+err.Error())
				log.Println(err)
			}

			var messages = chunk(msg, 2000)
			for _, v := range messages {
				msg, err := s.ChannelMessageSend(m.ChannelID, v)
				if err != nil {
					panic(err)
				}
				groqLUT[msg.ID] = history
			}
		case "groq-prompt":
			role := strings.Join(tokens[1:], " ")

			history := func() []Message {
				if m.ReferencedMessage != nil {
					var ctx = groqLUT[m.ReferencedMessage.ID]
					return AskGroqSystem(role, ctx)
				} else {
					return AskGroqSystem(role, nil)
				}
			}()
			msg, err := s.ChannelMessageSend(m.ChannelID, "My role is now:"+role)
			if err != nil {
				panic(err)
			}
			groqLUT[msg.ID] = history

		default:
		}
	}
}

func chunk(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}
