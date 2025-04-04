package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls/v2"
)

// Current song API response
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

const api = `https://r-a-d.io/api`
const frontpage = `https://r-a-d.io`
const radioRed = 0xDF4C3A

func radioState() (Current, error) {
	//Get the current state structure from the r/a/dio API
	resp, err := http.Get(api)
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
	queue := make([]*discordgo.MessageEmbedField, len(record.Main.Queue))
	for i := range queue {
		timeLeft := (time.Duration(time.Unix(record.Main.Queue[i].Timestamp, 0).Sub(time.Now().In(time.UTC)).Seconds()) * time.Second).String()
		for i := 0; i < len(timeLeft)-10; i++ {
			timeLeft = " " + timeLeft
		}
		queue[i] = new(discordgo.MessageEmbedField)
		queue[i].Name = timeLeft + " from now"
		queue[i].Value = record.Main.Queue[i].Meta
		queue[i].Inline = false
	}

	embed := &discordgo.MessageEmbed{URL: frontpage,
		Title:     "Playback queue",
		Color:     radioRed,
		Footer:    &discordgo.MessageEmbedFooter{Text: "Now playing: " + record.Main.Np},
		Author:    &discordgo.MessageEmbedAuthor{Name: record.Main.Djname, IconURL: api + "/dj-image/" + record.Main.Dj.Djimage},
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: api + "/dj-image/" + record.Main.Dj.Djimage},
		Fields:    queue,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func radioCurrent(s *discordgo.Session, m *discordgo.MessageCreate) {
	record, err := radioState()
	if err != nil {
		log.Println(err)
		return
	}

	progress := (time.Duration(record.Main.Current-record.Main.StartTime) * time.Second).String() + " / " + (time.Duration(record.Main.EndTime-record.Main.StartTime) * time.Second).String()

	fields := make([]*discordgo.MessageEmbedField, 2)
	fields[0] = new(discordgo.MessageEmbedField)
	fields[1] = new(discordgo.MessageEmbedField)

	fields[1].Name = record.Main.Dj.Djname
	fields[1].Inline = false
	fields[1].Value = "Listeners: " + strconv.Itoa(record.Main.Listeners)
	fields[0].Name = record.Main.Np
	fields[0].Value = progress
	fields[0].Inline = false

	if !record.Main.Isafkstream && !IsUrl(record.Main.Thread) && HasImage(record.Main.Thread) {

		embed := &discordgo.MessageEmbed{URL: frontpage,
			Title:     "Now playing",
			Color:     radioRed,
			Image:     GetImage(record.Main.Thread),
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: api + "/dj-image/" + record.Main.Dj.Djimage},
			Fields:    fields,
		}
		if _, err := s.ChannelMessageSendEmbed(m.ChannelID, embed); err != nil {
			log.Println(err)
		}

	} else {
		footer := new(discordgo.MessageEmbedFooter)
		if !record.Main.Isafkstream {
			footer.Text = "Current thread: " + record.Main.Thread
		} else {
			footer.Text = "Upcoming: " + record.Main.Queue[0].Meta
		}

		embed := &discordgo.MessageEmbed{URL: frontpage,
			Title:     "Now playing",
			Color:     radioRed,
			Footer:    footer,
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: api + "/dj-image/" + record.Main.Dj.Djimage},
			Fields:    fields,
		}
		if _, err := s.ChannelMessageSendEmbed(m.ChannelID, embed); err != nil {
			log.Println(err)
		}
	}

}

func GetImage(s string) *discordgo.MessageEmbedImage {
	rxRelaxed := xurls.Relaxed()
	res := rxRelaxed.FindString(s)

	return &discordgo.MessageEmbedImage{
		URL: res,
	}
}

func HasImage(s string) bool {
	rxRelaxed := xurls.Relaxed()
	res := rxRelaxed.FindString(s)
	if res == "" {
		return false
	}
	ext, err := GetFileExtensionFromUrl(res)
	if err != nil {
		return false
	}

	things := []string{"gif", "jpg", "png", "webp", "jpeg"}
	return slices.Contains(things, ext)

}

func GetFileExtensionFromUrl(rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	pos := strings.LastIndex(u.Path, ".")
	if pos == -1 {
		return "", errors.New("couldn't find a period to indicate a file extension")
	}
	return u.Path[pos+1 : len(u.Path)], nil
}

func IsUrl(s string) bool {
	_, err := url.Parse(s)
	return err != nil
}

// Songs is the API response of the current music queue
type Songs []struct {
	Artist        string `json:"artist"`
	Title         string `json:"title"`
	ID            int    `json:"id"`
	Lastplayed    int    `json:"lastplayed"`
	Lastrequested int    `json:"lastrequested"`
	Requestable   bool   `json:"requestable"`
}

// Search is the API response for a music index search query
type Search struct {
	Total       int   `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	From        int   `json:"from"`
	To          int   `json:"to"`
	Data        Songs `json:"data"`
}

func radioSearchResults(query string) Songs {
	const searchAPI = api + "/search/"
	//Get the current state structure from the r/a/dio API
	resp, err := http.Get(searchAPI + query)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Error: " + http.StatusText(resp.StatusCode))
	}
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record Search

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	var songs Songs
	songs = append(songs, record.Data...)

	for i := 2; i <= record.LastPage; i++ {
		resp, err := http.Get(searchAPI + query + "?page=" + strconv.Itoa(i))
		if err != nil {
			log.Println(err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Println("Error: " + http.StatusText(resp.StatusCode))
		}

		// Fill the record with the data from the JSON
		var record Search

		// Use json.Decode for reading streams of JSON data
		if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
			log.Println(err)
		}
		songs = append(songs, record.Data...)
		resp.Body.Close()
	}
	return songs
}

func radioSearch(s *discordgo.Session, m *discordgo.MessageCreate, name string) {
	songs := radioSearchResults(name)
	lines := make([]string, len(songs))
	for i := range lines {
		lines[i] = songs[i].Artist + " - " + "[" + songs[i].Title + "]" + "(" + `https://r-a-d.io/request/` + strconv.Itoa(songs[i].ID) + ")"
	}
	if _, err := s.ChannelMessageSend(m.ChannelID, strings.Join(lines, "\n")); err != nil {
		log.Println(err)
	}
}
