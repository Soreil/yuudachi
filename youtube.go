package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

//MaxYoutubeResults is a configurable limit to amount of videos to get from API
var MaxYoutubeResults = 1

func youtubeSearch(s *discordgo.Session, m *discordgo.MessageCreate, query string) {

	var ctx context.Context
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(*youtubeAPIKey))
	if err != nil {
		log.Println(err)
		return
	}
	call := youtubeService.Search.List("id,snippet").Q(query).MaxResults(int64(MaxYoutubeResults))
	call.Type("video")
	response, err := call.Do()
	if err != nil {
		log.Println(err)
		return
	}

	var msgs []string
	msgs = append(msgs, "Youtube search results:")

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			var msg string
			//msg += item.Snippet.ChannelTitle + "\n"
			//msg += item.Snippet.Title + "\n"
			msg += fmt.Sprintf("https://youtube.com/watch/?v=%s\n", item.Id.VideoId)

			msgs = append(msgs, msg)
		}
	}
	channelMessageSendDeleteAble(s, m, strings.Join(msgs, "\n"))
}
