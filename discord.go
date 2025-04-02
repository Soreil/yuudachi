package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/soreil/yuudachi/groq"
)

const prefix string = "!"
const defaultPrompt string = `Be succinct\n`

var groqLUT map[string][]groq.Message = map[string][]groq.Message{}
var messageModel map[string]groq.ReasoningFormats = map[string]groq.ReasoningFormats{}

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
		msgModel := messageModel[m.ReferencedMessage.ID]

		if contains {
			msg, history, err := groq.AskGroq(m.Content, ctx, msgModel)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Oopsie woopsie we got an error groq sisters:"+err.Error())
				log.Println(err)
			}

			if msgModel == groq.RawReasoningFormat {
				thinkclosetag := `</think>`
				idx := strings.Index(msg, `</think>`)
				thinkingBlock := msg[7:idx]
				bodyStart := idx + len(thinkclosetag)
				body := strings.TrimSpace(msg[bodyStart:])
				reader := strings.NewReader(strings.TrimSpace(thinkingBlock))

				var files []*discordgo.File = make([]*discordgo.File, 1)
				files[0] = &discordgo.File{Name: "thinking.txt", ContentType: "text/html", Reader: reader}

				var messages = chunk(body, 1000)
				for _, v := range messages {
					msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Content: v, Files: files})

					if err != nil {
						panic(err)
					}
					groqLUT[msg.ID] = history
				}
			} else {
				var messages = chunk(msg, 2000)
				for _, v := range messages {
					msg, err := s.ChannelMessageSend(m.ChannelID, v)

					if err != nil {
						panic(err)
					}
					groqLUT[msg.ID] = history
				}
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
		case "reason":
			if !usingDeepSeekR1 {
				break
			}
			body := strings.Join(tokens[1:], " ")
			if body == "" {
				return
			}
			msg, history, err := func() (string, []groq.Message, error) {
				if m.ReferencedMessage != nil {
					var ctx = groqLUT[m.ReferencedMessage.ID]
					return groq.AskGroq(body, ctx, groq.RawReasoningFormat)
				} else {
					if usingDeepSeekR1 {
						var message = defaultPrompt + body
						return groq.AskGroq(message, nil, groq.RawReasoningFormat)
					}
				}
				panic("eek")
			}()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Oopsie woopsie we got an error groq sisters:"+err.Error())
				log.Println(err)
			}

			thinkclosetag := `</think>`
			idx := strings.Index(msg, `</think>`)
			thinkingBlock := msg[7:idx]
			bodyStart := idx + len(thinkclosetag)
			responsebody := strings.TrimSpace(msg[bodyStart:])
			reader := strings.NewReader(strings.TrimSpace(thinkingBlock))

			var files []*discordgo.File = make([]*discordgo.File, 1)
			files[0] = &discordgo.File{Name: "thinking.txt", ContentType: "text/html", Reader: reader}

			var messages = chunk(responsebody, 2000)
			for _, v := range messages {
				msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Content: v, Files: files})

				if err != nil {
					panic(err)
				}
				groqLUT[msg.ID] = history
				messageModel[msg.ID] = groq.RawReasoningFormat
			}

		case "groq":
			body := strings.Join(tokens[1:], " ")
			if body == "" {
				return
			}
			msg, history, err := func() (string, []groq.Message, error) {
				if m.ReferencedMessage != nil {
					var ctx = groqLUT[m.ReferencedMessage.ID]
					return groq.AskGroq(body, ctx, groq.HiddenReasoningFormat)
				} else {
					if usingDeepSeekR1 {
						var message = defaultPrompt + body
						return groq.AskGroq(message, nil, groq.HiddenReasoningFormat)
					} else {
						var x = groq.AskGroqSystem(defaultPrompt, nil)
						return groq.AskGroq(body, x, groq.HiddenReasoningFormat)
					}
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
				messageModel[msg.ID] = groq.HiddenReasoningFormat
			}
		case "groq-prompt":
			role := strings.Join(tokens[1:], " ")

			history := func() []groq.Message {
				if m.ReferencedMessage != nil {
					var ctx = groqLUT[m.ReferencedMessage.ID]
					return groq.AskGroqSystem(role, ctx)
				} else {
					return groq.AskGroqSystem(role, nil)
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
