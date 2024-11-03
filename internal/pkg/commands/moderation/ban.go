package moderation

import (
	"fmt"
	"strings"

	"github.com/Paranoia8972/PixelBot/internal/pkg/commands"
	"github.com/bwmarrin/discordgo"
)

func BanCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		commands.RespondWithMessage(s, i, "Please provide a user to ban.")
		return
	}

	user := options[0].UserValue(s)
	reason := fmt.Sprintf("%s: No reason provided", i.Member.User.Username)

	if len(options) > 1 {
		reason = fmt.Sprintf("%s: %s", i.Member.User.Username, options[1].StringValue())
	}

	err := s.GuildBanCreateWithReason(i.GuildID, user.ID, reason, 0)
	if err != nil {
		commands.RespondWithMessage(s, i, "Failed to ban user.")
		return
	}

	commands.RespondWithMessage(s, i, fmt.Sprintf("Successfully banned %s#%s\nReason: %s", user.Username, user.Discriminator, reason))
}

func UnbanCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		commands.RespondWithMessage(s, i, "Please provide a user to unban.")
		return
	}

	userID := options[0].StringValue()

	err := s.GuildBanDelete(i.GuildID, userID)
	if err != nil {
		commands.RespondWithMessage(s, i, "Failed to unban user.")
		return
	}

	commands.RespondWithMessage(s, i, fmt.Sprintf("Successfully unbanned user with ID %s", userID))
}

func UnbanAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	query := i.ApplicationCommandData().Options[0].StringValue()
	bans, err := s.GuildBans(i.GuildID, 0, "", "")
	if err != nil {
		commands.RespondWithMessage(s, i, "Failed to fetch banned users.")
		return
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, ban := range bans {
		if strings.Contains(strings.ToLower(ban.User.Username), strings.ToLower(query)) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  fmt.Sprintf("%s#%s", ban.User.Username, ban.User.Discriminator),
				Value: ban.User.ID,
			})
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		fmt.Printf("Failed to respond to autocomplete interaction: %v\n", err)
	}
}
