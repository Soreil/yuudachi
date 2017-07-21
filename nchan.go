package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jaytaylor/html2text"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"unicode"
)

// Adds <> around the links to prevent embedding
//TODO(sjon): Cleanup and verify correctness
func unembedURL(s string) string {
	start := strings.Index(s, "http")
	if start == -1 {
		return s
	}
	end := strings.IndexFunc(s[start:], unicode.IsSpace) + start
	if end > start {
		return s[:start] + "<" + s[start:end] + ">" + unembedURL(s[end:])
	}
	return s
}

func fourchan(s *discordgo.Session, m *discordgo.MessageCreate, board string) {
	type Catalog []struct {
		Page    int `json:"page"`
		Threads []struct {
			No            int    `json:"no"`
			Sticky        int    `json:"sticky,omitempty"`
			Closed        int    `json:"closed,omitempty"`
			Now           string `json:"now"`
			Name          string `json:"name"`
			Sub           string `json:"sub,omitempty"`
			Com           string `json:"com"`
			Filename      string `json:"filename"`
			Ext           string `json:"ext"`
			W             int    `json:"w"`
			H             int    `json:"h"`
			TnW           int    `json:"tn_w"`
			TnH           int    `json:"tn_h"`
			Tim           int64  `json:"tim"`
			Time          int    `json:"time"`
			Md5           string `json:"md5"`
			Fsize         int    `json:"fsize"`
			Resto         int    `json:"resto"`
			ID            string `json:"id"`
			Country       string `json:"country"`
			SemanticURL   string `json:"semantic_url"`
			CountryName   string `json:"country_name"`
			Replies       int    `json:"replies"`
			Images        int    `json:"images"`
			LastModified  int    `json:"last_modified"`
			Bumplimit     int    `json:"bumplimit,omitempty"`
			Imagelimit    int    `json:"imagelimit,omitempty"`
			OmittedPosts  int    `json:"omitted_posts,omitempty"`
			OmittedImages int    `json:"omitted_images,omitempty"`
			LastReplies   []struct {
				No          int    `json:"no"`
				Now         string `json:"now"`
				Name        string `json:"name"`
				Com         string `json:"com"`
				Time        int    `json:"time"`
				Resto       int    `json:"resto"`
				ID          string `json:"id"`
				Country     string `json:"country"`
				CountryName string `json:"country_name"`
				Filename    string `json:"filename,omitempty"`
				Ext         string `json:"ext,omitempty"`
				W           int    `json:"w,omitempty"`
				H           int    `json:"h,omitempty"`
				TnW         int    `json:"tn_w,omitempty"`
				TnH         int    `json:"tn_h,omitempty"`
				Tim         int64  `json:"tim,omitempty"`
				Md5         string `json:"md5,omitempty"`
				Fsize       int    `json:"fsize,omitempty"`
			} `json:"last_replies,omitempty"`
			MImg int `json:"m_img,omitempty"`
		} `json:"threads"`
	}
	const (
		apiRoot     = "https://a.4cdn.org"
		usage       = "4chan usage:\n!4chan BOARD\nWhere category is one of the standard boards."
		catalogRoot = "catalog.json"
	)

	//TODO(sjon): Implement this in a cleaner manner.
	if board == "help" {
		ChannelMessageSendDeleteAble(s, m, usage)
		return
	}

	boardCatalog := apiRoot + "/" + board + "/" + catalogRoot
	resp, err := http.Get(boardCatalog)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	if err == nil && resp.StatusCode != http.StatusOK {
		log.Println("Error: " + http.StatusText(resp.StatusCode))
		return
	}
	defer resp.Body.Close()

	var record Catalog
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	//Extract random board page
	page := record[rand.Intn(len(record)-1)]
	//Extract random thread from the page
	thread := page.Threads[rand.Intn(len(page.Threads)-1)]
	//Try and filter out general threads, this method is awfully poor
	for strings.Contains(strings.ToLower(thread.Sub), "general") {
		page = record[rand.Intn(len(record)-1)]
		thread = page.Threads[rand.Intn(len(page.Threads)-1)]
	}

	//Clean up text from the selected thread
	thread.Sub, err = html2text.FromString(strings.Replace(thread.Sub, "<wbr>", "", -1))
	if err != nil {
		log.Println(err)
		return
	}

	thread.Com, err = html2text.FromString(strings.Replace(thread.Com, "<wbr>", "", -1))
	if err != nil {
		log.Println(err)
		return
	}

	//Discord formatting for thread titles
	if thread.Sub != "" {
		thread.Sub = "*" + thread.Sub + "*" + "\n"
	}

	//Formatting URLs for Discord
	thread.Com = unembedURL(thread.Com)

	//We currently don't fetch images to upload to Discord ourselves.
	img := fmt.Sprintf("https://i.4cdn.org/%s/%d%s", board, thread.Tim, thread.Ext)
	//Thread link formatted for Discord
	link := fmt.Sprintf("<https://i.4cdn.org/%s/thread/%d>", board, thread.No)
	//To be replaced by a Discord Embed structure
	ChannelMessageSendDeleteAble(s, m, fmt.Sprintf("%s\n%s\n\n%s", thread.Sub, thread.Com, link))

	//The file should be part of the embed struct too
	imgresp, err := http.Get(img)

	if err != nil {
		log.Println(err)
		return
	}
	defer imgresp.Body.Close()

	if imgresp.StatusCode != http.StatusOK {
		log.Println("Failed to fetch: " + http.StatusText(imgresp.StatusCode))
		return
	}

	s.ChannelFileSend(m.ChannelID, thread.Filename+thread.Ext, imgresp.Body)
}