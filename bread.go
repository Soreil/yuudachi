package main

import (
	"net/http"
	"log"
	"errors"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"math"
)

const breadBase = "RON"
const breadRatio = 1.0

const currencyAPI = `https://api.fixer.io/latest?base=` + breadBase

type Currency struct {
	Base  string `json:"base"`
	Date  string `json:"date"`
	Rates map[string]interface{} `json:"rates"`
}

func inBreads(amount float64, currency string) (float64, error) {
	resp, err := http.Get(currencyAPI)
	if err != nil {
		log.Println(err)
		return 0, errors.New("Failed to get bread API")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(http.StatusText(resp.StatusCode))
	}
	var rates Currency
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		panic(err)
	}

	if r, ok := rates.Rates[currency].(float64); !ok {
		return r, errors.New(currency + ": Unknown currency")
	}
	return breadRatio / rates.Rates[currency].(float64) * amount, nil
}

func breads(s *discordgo.Session, m *discordgo.MessageCreate, amount float64, currency string) {
	b, err := inBreads(amount, currency)
	if err != nil {
		ChannelMessageSendDeleteAble(s, m, err.Error())
	} else {
		ChannelMessageSendDeleteAble(s, m, "That's "+strconv.Itoa(int(b))+" breads and "+strconv.Itoa(int((b-math.Trunc(b))*20))+" slices.")
	}
}
