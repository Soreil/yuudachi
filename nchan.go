package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jaytaylor/html2text"
)

const (
	yotsubaBlue  = 0xEEF2FF
	yotsubaRed   = 0xFFFFEE
	yotsubaGreen = 0x35B214

	gifBanners = 253
	pngBanners = 262
	jpgBanners = 224
)

const (
	apiRoot     = "https://a.4cdn.org"
	chanUsage   = "4chan usage:\n!4chan BOARD\nWhere BOARD is one of the standard board names."
	catalogRoot = "catalog.json"
	bannerRoot  = "http://s.4cdn.org/image/title"
	countryRoot = "http://s.4cdn.org/image/country"
)

var bannedBoards = []string{"cm", "y", "gif", "e", "h", "hc", "b", "mlp", "lgbt", "soc", "s", "hm", "d", "t", "aco", "r", "pol", "trash"}

func BannedBoard(s string) bool {
	for _, v := range bannedBoards {
		if v == s {
			return true
		}
	}
	return false
}

type thread struct {
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
}

func fourchan(s *discordgo.Session, m *discordgo.MessageCreate, board string) {
	type Catalog []struct {
		Page    int      `json:"page"`
		Threads []thread `json:"threads"`
	}

	//TODO(sjon): Implement this in a cleaner manner.
	if board == "help" {
		s.ChannelMessageSend(m.ChannelID, chanUsage)
		return
	}

	boardCatalog := apiRoot + "/" + board + "/" + catalogRoot
	resp, err := http.Get(boardCatalog)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
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
	for strings.Contains(strings.ToLower(thread.Sub), "general") || thread.Sticky == 1 && thread.Replies < 10 {
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

	reply := formatThread(thread, board)
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, reply); err != nil {
		log.Println(err)
	}
}

func formatThread(thread thread, board string) *discordgo.MessageEmbed {
	img := fmt.Sprintf("https://i.4cdn.org/%s/%d%s", board, thread.Tim, thread.Ext)
	//thumb := fmt.Sprintf("https://i.4cdn.org/%s/%ds%s", board, thread.Tim, ".jpg")
	link := fmt.Sprintf("https://i.4cdn.org/%s/thread/%d", board, thread.No)

	var banner string
	switch rand.Intn(2) {
	case 0:
		banner = fmt.Sprintf(bannerRoot+"/%d%s", rand.Intn(jpgBanners), ".jpg")
	case 1:
		banner = fmt.Sprintf(bannerRoot+"/%d%s", rand.Intn(pngBanners), ".png")
	case 2:
		banner = fmt.Sprintf(bannerRoot+"/%d%s", rand.Intn(gifBanners), ".gif")
	}

	var replies string
	if thread.Replies == 0 {
		replies = "There are no replies ;_;"
	} else if thread.Replies == 1 {
		replies = "There is one reply."
	} else {
		replies = fmt.Sprintf("There are %d replies.", thread.Replies)
	}

	var title string

	if thread.Sub == "" {
		title = "link"
	} else {
		title = thread.Sub
	}

	embed := &discordgo.MessageEmbed{URL: link,
		Title:       title,
		Color:       yotsubaGreen,
		Footer:      &discordgo.MessageEmbedFooter{Text: replies},
		Author:      &discordgo.MessageEmbedAuthor{Name: thread.Name, IconURL: countryRoot + "/" + strings.ToLower(thread.Country) + ".gif"},
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: banner},
		Description: thread.Com,
		Image:       &discordgo.MessageEmbedImage{URL: img, Width: thread.W, Height: thread.H}}

	return embed
}
