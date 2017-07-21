package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os/exec"
	"strings"
)

const appVersion = `21-7-2017
*4chan embeds*`
const useFiglet = false

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
	if useFiglet {
		fig, err := figlet(msg)
		if err != nil {
			log.Println("Failed to create figlet text:", err)
		} else {
			ChannelMessageSendDeleteAble(s, m, "```"+fig+"```")
			return
		}
	}
	ChannelMessageSendDeleteAble(s, m, msg)
}

func usage(s *discordgo.Session, m *discordgo.MessageCreate) {
	usage := strings.Join([]string{"twitter", "version", "fortune" /*"8chan",*/, "4chan", "bible", "radio", "bird"}, ", ")
	ChannelMessageSendDeleteAble(s, m, "The possible commands Yuudachi will like: "+usage+".")
}
