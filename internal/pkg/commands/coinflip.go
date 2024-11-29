package commands

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func CoinFlipCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	outcome := "heads"
	if r.Intn(2) == 0 {
		outcome = "tails"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "The coin landed on " + outcome,
		},
	})
}
