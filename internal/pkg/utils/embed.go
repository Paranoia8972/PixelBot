package utils

import "github.com/bwmarrin/discordgo"

type EmbedType string

const (
	EmbedTypeWarn    EmbedType = "warn"
	EmbedTypeInfo    EmbedType = "info"
	EmbedTypeDefault EmbedType = "default"
)

// EmbedColor maps embed types to colors
var EmbedColor = map[EmbedType]int{
	EmbedTypeWarn:    0xff0000, // Red
	EmbedTypeInfo:    0x0000ff, // Blue
	EmbedTypeDefault: 0x00ff00, // Green
}

func CreateEmbed(title, description string, embedType EmbedType, thumbnailURL string) *discordgo.MessageEmbed {
	color, exists := EmbedColor[embedType]
	if !exists {
		color = EmbedColor[EmbedTypeDefault]
	}

	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: thumbnailURL,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cloud.onthepixel.net/apps/files_sharing/publicpreview/skQoiXkskNFyBBx?file=/&fileId=2125&x=1920&y=1080&a=true&etag=2ea48a07a9e7f79eadbaadfe61c44f75",
		},
	}
}
