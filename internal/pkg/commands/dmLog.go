package commands

import (
	"github.com/Paranoia8972/PixelBot/internal/pkg/utils"
	"github.com/bwmarrin/discordgo"
)

func DMLogCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel := i.ApplicationCommandData().Options[0].ChannelValue(s)
	guildID := i.GuildID

	err := utils.SetDMLogChannel(guildID, channel.ID)
	if err != nil {
		RespondWithMessage(s, i, "Failed to set DM log channel.")
		return
	}

	RespondWithMessage(s, i, "DM log channel set successfully!")
}
