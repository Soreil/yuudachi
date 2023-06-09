package main

import (
	"encoding/json"
	"log"
	"net/http"

	"errors"

	"github.com/bwmarrin/discordgo"
)

type Fortune struct {
	Fortune string `json:"fortune"`
}

type FortuneBodyError struct {
	Err error
}

func (f *FortuneBodyError) Error() string {
	return f.Err.Error()
}

func fortune(s *discordgo.Session, m *discordgo.MessageCreate, category string) {
	u := "http://www.yerkee.com/api/fortune"
	if category == "help" {
		s.ChannelMessageSend(m.ChannelID, "Fortune usage:\n!fortune CATEGORY\nWhere category is one of: computers, cookie, definitions, miscellaneous, people, platitudes, politics, science, wisdom\nIf no category is given a random one is chosen.")
		return
	}
	switch category {
	case "computers", "cookie", "definitions", "miscellaneous", "people", "platitudes", "politics", "science", "wisdom":
		u += "/" + category
	case "":
	default:
		s.ChannelMessageSend(m.ChannelID, "Unknown category, type \"!fortune help\" for a list of categories allowed.")
		return
	}

	record, err := FetchFortune(u)

	//Didn't get a fortune
	if err == nil {
		s.ChannelMessageSend(m.ChannelID, record.Fortune)
	} else {
		log.Println(err)
	}
}

func FetchFortune(url string) (Fortune, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Fortune{}, err
	}
	defer resp.Body.Close()

	var record Fortune
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	if record.Fortune == "" {
		return Fortune{}, errors.New("no fortune body")
	}

	return record, nil
}
