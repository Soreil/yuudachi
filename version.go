package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os/exec"
	"strings"
)

const appVersion = `26-4-2017
Delet this`

func figlet(s string) (string, error) {
	cmd := exec.Command("figlet", "-p", s)
	fig, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(fig), nil
}

func version(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "Current version: " + appVersion
	log.Println(msg)
	fig, err := figlet(msg)
	if err != nil {
		ChannelMessageSendDeleteAble(s, m, msg)
	} else {
		ChannelMessageSendDeleteAble(s, m, "```"+fig+"```")
	}
}

func usage(s *discordgo.Session, m *discordgo.MessageCreate) {
	usage := strings.Join([]string{"twitter", "version", "fortune" /*"8chan",*/, "4chan", "bible", "radio", "bird"}, ", ")
	ChannelMessageSendDeleteAble(s, m, "The possible commands Yuudachi will like: "+usage+".")
}
