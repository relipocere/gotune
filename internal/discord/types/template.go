package types

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	thumbnailURL     = "https://cdn.discordapp.com/attachments/902158239825788999/902159912967241738/gopher3d.png"
	colorGo      int = 1500402
	colorRed     int = 15214375
)

//TextInteractionResp ...
func TextInteractionResp(message string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}
}

//EmbedInteractionResp ...
func EmbedInteractionResp(embed *discordgo.MessageEmbed) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		},
	}
}

//QueueEmbed ...
func QueueEmbed(songs []Song) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       "Queue",
		Description: "The queue is empty",
		Color:       colorGo,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: thumbnailURL,
		},
	}

	if songs != nil && len(songs) > 0 {
		var list string
		for n, song := range songs {
			list += fmt.Sprintf("%d. %s\n", n+1, song.Title)
		}
		embed.Description = list
	}
	return embed
}

//TrackEmbed ...
func TrackEmbed(message string, s Song) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		URL:   s.Link,
		Title: s.Title,
		Color: colorGo,
		Author: &discordgo.MessageEmbedAuthor{
			Name: message,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: thumbnailURL,
		},
	}

	if s.Requester != nil {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Requested by %s", s.Requester.Username),
			IconURL: s.Requester.AvatarURL("32x32"),
		}
	}
	return embed
}

//ErrorEmbed ...
func ErrorEmbed(message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Error",
		Description: message,
		Color:       colorRed,
	}
}
