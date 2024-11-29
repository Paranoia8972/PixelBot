package commands

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ChooserCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a list of usernames.",
			},
		})
		return
	}

	usernames := strings.Split(options[0].StringValue(), ",")
	for i := range usernames {
		usernames[i] = strings.TrimSpace(usernames[i])
	}

	if len(usernames) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No valid usernames provided.",
			},
		})
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	chosen := usernames[r.Intn(len(usernames))]

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Chosen username: " + chosen,
		},
	})
}
