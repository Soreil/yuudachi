package main

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func roll(s *discordgo.Session, m *discordgo.MessageCreate, n int) {
	val := rand.Intn(n) + 1
	channelMessageSendDeleteAble(s, m, fmt.Sprint(val))
}
