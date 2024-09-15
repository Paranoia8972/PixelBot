package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/paranoia8972/PixelBot/utils"
)

func AutoRoleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if len(i.ApplicationCommandData().Options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No subcommand provided.",
			},
		})
		return
	}

	switch i.ApplicationCommandData().Options[0].Name {
	case "add":
		AutoRoleAddCommand(s, i)
	case "get":
		AutoRoleGetCommand(s, i)
	case "remove":
		AutoRoleRemoveCommand(s, i)
	default:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unknown subcommand.",
			},
		})
	}
}

func AutoRoleAddCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a role to add.",
			},
		})
		return
	}

	roleID := options[0].RoleValue(s, "").ID

	err := utils.AddAutoRole(i.GuildID, roleID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to add auto role.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Auto role added successfully!",
		},
	})
}

func AutoRoleGetCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	autoRoles, err := utils.GetAutoRoles(i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve auto roles.",
			},
		})
		return
	}

	var roleMentions []string
	for _, roleID := range autoRoles {
		roleMentions = append(roleMentions, "<@&"+roleID+">")
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Auto Roles",
					Description: "Auto roles for this server.",
					Color:       0x00ff00,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Roles",
							Value:  "• " + strings.Join(roleMentions, "\n• "),
							Inline: true,
						},
					},
				},
			},
		},
	})
}

func AutoRoleRemoveCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a role to remove.",
			},
		})
		return
	}

	roleID := options[0].RoleValue(s, "").ID

	err := utils.RemoveAutoRole(i.GuildID, roleID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to remove auto role.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Auto role removed successfully!",
		},
	})
}
