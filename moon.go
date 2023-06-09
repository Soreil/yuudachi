package main

import (
	"math"
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

func MakeMoon() string {
	previousNewMoon := time.Date(2019, time.June, 3, 12, 1, 0, 0, time.Local)
	delta := time.Since(previousNewMoon)
	daysPerMonth := 29.530588853
	moonDay := math.Mod(delta.Hours()/24, daysPerMonth) + 1
	percentage := moonDay / daysPerMonth
	return string(percentPhase(percentage))
}

func moonPhase(s *discordgo.Session, m *discordgo.MessageCreate) {

	moon := MakeMoon()

	s.ChannelMessageSend(m.ChannelID, moon)
}
