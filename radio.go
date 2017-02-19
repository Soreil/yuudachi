package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"strings"
	"time"
)

func radio(s *discordgo.Session, m *discordgo.MessageCreate, function string) {
	type Current struct {
		Main struct {
			Np           string `json:"np"`
			Listeners    int    `json:"listeners"`
			Bitrate      int    `json:"bitrate"`
			Isafkstream  bool   `json:"isafkstream"`
			Isstreamdesk bool   `json:"isstreamdesk"`
			Current      int    `json:"current"`
			StartTime    int    `json:"start_time"`
			EndTime      int    `json:"end_time"`
			Lastset      string `json:"lastset"`
			Trackid      int    `json:"trackid"`
			Thread       string `json:"thread"`
			Requesting   bool   `json:"requesting"`
			Djname       string `json:"djname"`
			Dj           struct {
				ID       int    `json:"id"`
				Djname   string `json:"djname"`
				Djtext   string `json:"djtext"`
				Djimage  string `json:"djimage"`
				Djcolor  string `json:"djcolor"`
				Visible  bool   `json:"visible"`
				Priority int    `json:"priority"`
				CSS      string `json:"css"`
				ThemeID  int    `json:"theme_id"`
				Role     string `json:"role"`
			} `json:"dj"`
			Queue []struct {
				Meta      string `json:"meta"`
				Time      string `json:"time"`
				Type      int    `json:"type"`
				Timestamp int    `json:"timestamp"`
			} `json:"queue"`
			Lp []struct {
				Meta      string `json:"meta"`
				Time      string `json:"time"`
				Type      int    `json:"type"`
				Timestamp int    `json:"timestamp"`
			} `json:"lp"`
		} `json:"main"`
	}

	const u = `https://r-a-d.io/api`

	// Ignore all messages created by the bot itself
	if m.Author.Bot || m.Author.Username == "Liru" {
		return
	}

	if function == "" {
		function = "dj"
	}

	if function == "help" {
		s.ChannelMessageSend(m.ChannelID, "Usage: !radio dj")
		return
	}
	resp, err := http.Get(u)
	if err != nil {
		log.Fatal("Get: ", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Error: " + http.StatusText(resp.StatusCode))
		return
	}
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record Current

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	var parts []string
	time.Duration(record.Main.Current - record.Main.StartTime).String()
	parts = append(parts, "Present DJ: "+record.Main.Dj.Djname)
	parts = append(parts, "Present song: "+record.Main.Np)
	parts = append(parts, "Present time: "+(time.Duration(record.Main.Current-record.Main.StartTime)*time.Second).String()+" / "+(time.Duration(record.Main.EndTime-record.Main.StartTime)*time.Second).String())
	if record.Main.Thread != "" && record.Main.Thread != "none" {
		parts = append(parts, "Present place: "+record.Main.Thread)
	} else {
		parts = append(parts, "There is no thread up at the moment.")
	}
	s.ChannelMessageSend(m.ChannelID, strings.Join(parts, "\n"))

	imgresp, err := http.Get(u + "/dj-image/" + record.Main.Dj.Djimage)

	if err != nil {
		log.Println(err)
		return
	}
	defer imgresp.Body.Close()

	if imgresp.StatusCode != http.StatusOK {
		log.Println("Failed to fetch: " + http.StatusText(imgresp.StatusCode))
		return
	}

	var format string
	switch t := imgresp.Header.Get("Content-Type"); {
	case strings.Contains(t, "png"):
		format = "png"
	case strings.Contains(t, "jpg"), strings.Contains(t, "jpeg"):
		format = "jpeg"
	case strings.Contains(t, "gif"):
		format = "gif"
	default:
		format = "unknown"
	}
	s.ChannelFileSend(m.ChannelID, record.Main.Dj.Djimage+"."+format, imgresp.Body)
}
