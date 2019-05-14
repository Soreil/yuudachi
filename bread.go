package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"math"
	"net/http"
	"strconv"
)
var fixerAPIToken string
const breadBase = "RON"
const breadRatio = 1.0

var currencyAPI = `https://data.fixer.io/api/latest?access_key=`+ fixerAPIToken + `?base=` + breadBase

type Currency struct {
	Base  string                 `json:"base"`
	Date  string                 `json:"date"`
	Rates map[string]interface{} `json:"rates"`
}

func rates() (Currency, error) {
	var rates Currency
	resp, err := http.Get(currencyAPI)
	if err != nil {
		log.Println(err)
		return rates, errors.New("Failed to get bread API")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return rates, errors.New(http.StatusText(resp.StatusCode))
	}
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		panic(err)
	}
	return rates, nil
}
func inBreads(amount float64, currency string) (float64, error) {
	rates, err := rates()
	if err != nil {
		return 0, err
	}
	if r, ok := rates.Rates[currency].(float64); !ok {
		return r, errors.New(currency + ": Unknown currency")
	}
	return breadRatio / rates.Rates[currency].(float64) * amount, nil
}

func fromBreads(amount float64, currency string) (float64, error) {
	rates, err := rates()
	if err != nil {
		return 0, err
	}
	if r, ok := rates.Rates[currency].(float64); !ok {
		return r, errors.New(currency + ": Unknown currency")
	}
	return breadRatio * rates.Rates[currency].(float64) * amount, nil
}

func breads(s *discordgo.Session, m *discordgo.MessageCreate, amount float64, currency string) {
	b, err := inBreads(amount, currency)
	if err != nil {
		ChannelMessageSendDeleteAble(s, m, err.Error())
	} else {
		ChannelMessageSendDeleteAble(s, m, "That's "+strconv.Itoa(int(b))+" breads and "+strconv.Itoa(int((b-math.Trunc(b))*20))+" slices.")
	}
}

func fiats(s *discordgo.Session, m *discordgo.MessageCreate, amount float64, currency string) {
	b, err := fromBreads(amount, currency)
	if err != nil {
		ChannelMessageSendDeleteAble(s, m, err.Error())
	} else {
		ChannelMessageSendDeleteAble(s, m, fmt.Sprintf("That's %02.2f %s.", b, currency))
	}
}
