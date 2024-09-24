package commands

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func PingCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	start := time.Now()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pinging...",
		},
	})
	if err != nil {
		return
	}

	latency := time.Since(start).Milliseconds()
	latencyMessage := "Pong! The bot's latency is " + strconv.FormatInt(latency, 10) + "ms"

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &latencyMessage,
	})
	if err != nil {
		return
	}
}
