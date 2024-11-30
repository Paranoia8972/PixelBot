package commands

import (
	"fmt"

	"github.com/Paranoia8972/PixelBot/internal/pkg/utils"
	"github.com/bwmarrin/discordgo"
)

func LevelCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	guildID := i.GuildID

	level, xp, err := utils.GetUserLevelAndXP(guildID, userID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error retrieving your level.",
				Flags:   64,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You are level %d with %d XP!", level, xp),
		},
	})
}
