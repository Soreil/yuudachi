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
		s.ChannelMessageSend(m.ChannelID, "Fortune usage:\n!fortune CATEGORY\nWhere category is one of: computers, cookie, definitions, miscellaneous, people, platitudes, politics, science, wisdom")
		return
	}

	if category == "computers" || category == "cookie" || category == "definitions" || category == "miscellaneous" || category == "people" || category == "platitudes" || category == "politics" || category == "science" || category == "wisdom" {
		u += "/" + category
	} else if category != "" {
		s.ChannelMessageSend(m.ChannelID, "Unknown category, type \"!fortune help\" for a list of categories allowed.")
		return
	}
	resp, err := http.Get(u)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record Fortune

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	if record.Fortune != "" {
		s.ChannelMessageSend(m.ChannelID, record.Fortune)
	}
}
