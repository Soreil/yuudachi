package main

import (
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"strings"
	"time"
)

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
			Timestamp int64  `json:"timestamp"`
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

func radioState() (Current, error) {
	//Get the current state structure from the r/a/dio API
	resp, err := http.Get(u)
	if err != nil {
		return Current{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Current{}, errors.New("Error: " + http.StatusText(resp.StatusCode))
	}
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record Current

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		return Current{}, err
	}
	return record, nil
}

func radioHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Usage: !radio [queue, dj]")
}

func radioQueue(s *discordgo.Session, m *discordgo.MessageCreate) {
	record, err := radioState()
	if err != nil {
		log.Println(err)
		return
	}
	if !record.Main.Isafkstream {
		s.ChannelMessageSend(m.ChannelID, "Sadly Hanyuu is not in a position to play the current queue, a DJ is playing!")
		return
	}
	queue := make([]string, len(record.Main.Queue)+1)
	queue[0] = "r/a/dio playback queue:"
	for i := range queue[1:] {
		queue[i+1] = (time.Duration(time.Unix(record.Main.Queue[i].Timestamp, 0).Sub(time.Now().In(time.UTC)).Seconds()) * time.Second).String() + ": " + record.Main.Queue[i].Meta

	}
	s.ChannelMessageSend(m.ChannelID, strings.Join(queue, "\n"))
}

func radioCurrent(s *discordgo.Session, m *discordgo.MessageCreate) {
	record, err := radioState()
	if err != nil {
		log.Println(err)
		return
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
