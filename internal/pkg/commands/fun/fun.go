package fun

import (
	"github.com/Paranoia8972/PixelBot/internal/app/config"
	"github.com/bwmarrin/discordgo"
)

var cfg *config.Config

func init() {
	cfg = config.LoadConfig()
}

func RespondWithMessage(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   64,
		},
	})
}
