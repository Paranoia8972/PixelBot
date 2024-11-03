package moderation

import (
	"fmt"

	"github.com/Paranoia8972/PixelBot/internal/pkg/commands"
	"github.com/bwmarrin/discordgo"
)

func KickCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		commands.RespondWithMessage(s, i, "Please provide a user to kick.")
		return
	}

	user := options[0].UserValue(s)
	reason := fmt.Sprintf("%s: No reason provided", i.Member.User.Username)

	if len(options) > 1 {
		reason = fmt.Sprintf("%s: %s", i.Member.User.Username, options[1].StringValue())
	}

	err := s.GuildMemberDeleteWithReason(i.GuildID, user.ID, reason)
	if err != nil {
		commands.RespondWithMessage(s, i, "Failed to kick user.")
		return
	}

	commands.RespondWithMessage(s, i, fmt.Sprintf("Successfully kicked %s#%s\nReason: %s", user.Username, user.Discriminator, reason))
}
