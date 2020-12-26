package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

type moonPhaseName string

const (
	newMoon            moonPhaseName = "New Moon"
	waxingCrescentMoon moonPhaseName = "Waxing Crescent Moon"
	firstQuarter       moonPhaseName = "First Quarter"
	waxingGibbousMoon  moonPhaseName = "Waxing Gibbous Moon"
	fullMoon           moonPhaseName = "Full Moon"
	waningGibbousMoon  moonPhaseName = "Waning Gibbous Moon"
	lastQuarter        moonPhaseName = "Last Quarter"
	waningCrescentMoon moonPhaseName = "Waning Crescent Moon"
)

const (
	newMoonRune            = '🌑'
	waxingCrescentMoonRune = '🌒'
	firstQuarterRune       = '🌓'
	waxingGibbousMoonRune  = '🌔'
	fullMoonRune           = '🌕'
	waningGibbousMoonRune  = '🌖'
	lastQuarterMoonRune    = '🌗'
	waningCrescentMoonRune = '🌘'
)

func (p moonPhaseName) Rune() rune {
	switch p {
	case newMoon:
		return newMoonRune
	case waxingCrescentMoon:
		return waxingCrescentMoonRune
	case firstQuarter:
		return firstQuarterRune
	case waxingGibbousMoon:
		return waxingGibbousMoonRune
	case fullMoon:
		return fullMoonRune
	case waningGibbousMoon:
		return waningGibbousMoonRune
	case lastQuarter:
		return lastQuarterMoonRune
	case waningCrescentMoon:
		return waningCrescentMoonRune
	default:
		return ' '
	}
}

const lunarMonthLength = time.Hour * 24 * 28 //28 days is a wild guess, easily divided by 4

//MoonResponse maps the API response to a Go struct
type MoonResponse struct {
	Error       bool   `json:"error"`
	Apiversion  string `json:"apiversion"`
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	Day         int    `json:"day"`
	Numphases   int    `json:"numphases"`
	Datechanged bool   `json:"datechanged"`
	Phasedata   []struct {
		Phase string `json:"phase"`
		Date  string `json:"date"`
		Time  string `json:"time"`
	} `json:"phasedata"`
}

const americanTime = `2006 Jan 2 15:04`

func closestPhase(response MoonResponse) rune {
	var lastTime time.Time
	var lastPhase string
	for _, phase := range response.Phasedata {
		phaseTime, err := time.Parse(americanTime, phase.Date+" "+phase.Time)
		if err != nil {
			panic(err)
		}

		if phaseTime.Before(time.Now()) {
			lastTime = phaseTime
			lastPhase = phase.Phase
			continue
		}

		//sanity check
		if !phaseTime.After(time.Now()) {
			panic("Logic error in time function!")
		}

		timeSinceLastPhase := time.Now().Sub(lastTime)
		timeUntilNextPhase := phaseTime.Sub(time.Now())
		difference := phaseTime.Sub(lastTime)

		if timeSinceLastPhase < timeUntilNextPhase {
			//We are closer to the moon phase before the current date than after it.
			if timeSinceLastPhase.Nanoseconds() > difference.Nanoseconds()/2 {
				//We actually want the "half"-phase after lastPhase
				if p := moonPhaseName(lastPhase).Rune(); p == waningCrescentMoonRune { //handling wrap around
					return newMoonRune
				}
				return moonPhaseName(lastPhase).Rune() + 1
			}
			//lastPhase was correct
			return moonPhaseName(lastPhase).Rune()

		}
		if timeUntilNextPhase.Nanoseconds() < difference.Nanoseconds()/2 {
			//We actually want the "half"-phase before phaseTime
			if p := moonPhaseName(phase.Phase).Rune(); p == newMoonRune { //handling wrap around
				return waningCrescentMoonRune
			}
			return moonPhaseName(phase.Phase).Rune() - 1
		}
		//nextPhase was correct
		return moonPhaseName(phase.Phase).Rune()
	}
	return ' '
}

func getMoonPhase() (moon rune, err error) {
	const url = "http://api.usno.navy.mil/moon/phase"

	year, month, day := time.Now().Add(-(lunarMonthLength / 2)).Date()                    // 14 days ago
	options := fmt.Sprintf("?year=%d&month=%d&day=%d&nump=4&dst=false", year, month, day) //Temporary awful

	req, err := http.NewRequest("GET", url+options, nil) //TODO(sjon): Figure out what is wrong with setting the header options.
	if err != nil {
		return moon, err
	}

	//req.Header.Add("year", strconv.Itoa(year))
	//req.Header.Add("month", strconv.Itoa(int(month)))
	//req.Header.Add("day", strconv.Itoa(day))
	//req.Header.Add("nump", strconv.Itoa(4)) //Number of moon phases requested
	//req.Header.Add("dst", "false")

	clt := http.Client{
		Timeout: time.Second,
	}
	resp, err := clt.Do(req)
	if err != nil {
		log.Println(err)
		return moon, err
	}
	if resp.StatusCode != http.StatusOK {
		return moon, errors.New("Moon phase lookup failed: " + resp.Status)
	}

	data, _ := ioutil.ReadAll(resp.Body)
	var moonReply MoonResponse
	err = json.Unmarshal(data, &moonReply)
	if err != nil {
		return moon, errors.New("Failed to parse moon phase API response: " + err.Error())
	}

	if len(moonReply.Phasedata) == 0 {
		return moon, errors.New("didn't get phase data from API response")
	}
	moon = closestPhase(moonReply)
	return
}

func percentPhase(percentage float64) rune {
	percentage *= 8
	percentage = math.Round(percentage)
	switch percentage {
	case 0:
		return newMoon.Rune()
	case 1:
		return waxingCrescentMoon.Rune()
	case 2:
		return firstQuarter.Rune()
	case 3:
		return waxingGibbousMoon.Rune()
	case 4:
		return fullMoon.Rune()
	case 5:
		return waningGibbousMoon.Rune()
	case 6:
		return lastQuarter.Rune()
	case 7:
		return waningCrescentMoon.Rune()
	case 8:
		return newMoon.Rune()
	default:
		return '🦆'

	}
}

func moonPhase(s *discordgo.Session, m *discordgo.MessageCreate) {
	moon, err := getMoonPhase()
	if err != nil {
		//log.Println(err)
		//channelMessageSendDeleteAble(s, m, "Sorry, I failed to inspect the moon.")
		//return
		previousNewMoon := time.Date(2019, time.June, 3, 12, 1, 0, 0, time.Local)
		delta := time.Since(previousNewMoon)
		daysPerMonth := 29.530588853
		moonDay := math.Mod(delta.Hours()/24, daysPerMonth) + 1
		percentage := moonDay / daysPerMonth
		moon = percentPhase(percentage)

	}
	channelMessageSendDeleteAble(s, m, string(moon))
}
