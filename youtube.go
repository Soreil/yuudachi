package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// MaxYoutubeResults is a configurable limit to amount of videos to get from API
var MaxYoutubeResults = 1

type videoQueue struct {
	userVideos map[string]<-chan string
}

var videos struct {
	sync.Mutex
	videoQueue
}

func init() {
	videos.videoQueue.userVideos = make(map[string]<-chan string)
}

func youtubeSearch(s *discordgo.Session, m *discordgo.MessageCreate, query string) {

	var ctx context.Context
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(*youtubeAPIKey))
	if err != nil {
		log.Println(err)
		return
	}
	call := youtubeService.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(int64(25))
	call.Type("video")
	response, err := call.Do()
	if err != nil {
		log.Println(err)
		return
	}

	var msgs []string

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			var msg string

			//Discord will create an embed structure from the link
			//so we don't need to provide our own metadata for video
			msg += fmt.Sprintf("https://youtube.com/watch/?v=%s\n", item.Id.VideoId)

			msgs = append(msgs, msg)
		}
	}
	for i := 0; i < MaxYoutubeResults && i < len(msgs); i++ {
		s.ChannelMessageSend(m.ChannelID, msgs[i])
	}

	if len(msgs) > 0 {
		vids := enqueueVideos(msgs[MaxYoutubeResults:])
		videos.Lock()
		defer videos.Unlock()
		videos.userVideos[m.Author.ID] = vids
	}
}

func enqueueVideos(videos []string) <-chan string {
	vidChan := make(chan string, len(videos))
	for _, v := range videos {
		vidChan <- v
	}
	return vidChan
}

func nextVideo(s *discordgo.Session, m *discordgo.MessageCreate) {
	videos.Lock()
	defer videos.Unlock()
	if len(videos.userVideos[m.Author.ID]) <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Your video queue is empty :(")
		return
	}
	s.ChannelMessageSend(m.ChannelID, <-videos.userVideos[m.Author.ID])
}
