package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func randomDogImage(s *discordgo.Session, m *discordgo.MessageCreate) {
	const dogAPIRoot = `https://dog.ceo/api/`
	const dogRandom = `breeds/image/random`

	type RandomDogResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	resp, err := http.Get(dogAPIRoot + dogRandom)
	if err != nil {
		log.Fatalln("Failed to fetch dog image:", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to fetch dog image: ", resp.Status)
	}

	defer resp.Body.Close()

	var dog RandomDogResponse
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &dog)

	log.Println(dog.Message)
	pic, err := http.Get(dog.Message)
	suffix := dog.Message[strings.LastIndex(dog.Message, "."):]
	s.ChannelFileSend(m.ChannelID, "Happy doggo"+suffix, pic.Body)
}
