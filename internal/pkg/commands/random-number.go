package commands

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func RandomNumberCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a number.",
			},
		})
		return
	}

	max := options[0].IntValue()
	if max < 1 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid number provided.",
			},
		})
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := r.Intn(int(max)) + 1

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Your random number is " + strconv.Itoa(randomNumber),
		},
	})
}
