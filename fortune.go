package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
)

func fortune(s *discordgo.Session, m *discordgo.MessageCreate, category string) {
	type Fortune struct {
		Fortune string `json:"fortune"`
	}
	u := "http://www.yerkee.com/api/fortune"
	if category == "help" {
		ChannelMessageSendDeleteAble(s, m, "Fortune usage:\n!fortune CATEGORY\nWhere category is one of: computers, cookie, definitions, miscellaneous, people, platitudes, politics, science, wisdom\nIf no category is given a random one is chosen.")
		return
	}
	switch category {
	case "computers", "cookie", "definitions", "miscellaneous", "people", "platitudes", "politics", "science", "wisdom":
		u += "/" + category
	case "":
	default:
		ChannelMessageSendDeleteAble(s, m, "Unknown category, type \"!fortune help\" for a list of categories allowed.")
		return
	}
	resp, err := http.Get(u)
	if err != nil {
		log.Fatal("Failed to get resource: ", err)
		return
	}
	defer resp.Body.Close()

	var record Fortune
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	//Didn't get a fortune
	if record.Fortune != "" {
		ChannelMessageSendDeleteAble(s, m, record.Fortune)
	} else {
		log.Println("Failed to get a fortune.")
	}
}
