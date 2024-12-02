package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SayHello(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User
	response := fmt.Sprintf("Hello, %s!", user.Mention())

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
