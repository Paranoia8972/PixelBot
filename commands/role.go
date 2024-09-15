package commands

import (
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func RoleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if len(i.ApplicationCommandData().Options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No subcommand provided.",
				Flags:   64,
			},
		})
		return
	}

	switch i.ApplicationCommandData().Options[0].Name {
	case "add":
		addRole(s, i)
	case "remove":
		removeRole(s, i)
	case "all":
		roleAll(s, i)
	}
}

func addRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var roleID, userID string

	for _, option := range i.ApplicationCommandData().Options[0].Options {
		switch option.Name {
		case "role":
			if option.Type != discordgo.ApplicationCommandOptionRole {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Invalid option type for role. Please provide a role.",
					},
				})
				return
			}
			roleID = option.RoleValue(s, "").ID
		case "user":
			if option.Type != discordgo.ApplicationCommandOptionUser {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Invalid option type for user. Please provide a user.",
					},
				})
				return
			}
			userID = option.UserValue(s).ID
		}
	}

	if roleID == "" || userID == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide both a role and a user.",
			},
		})
		return
	}

	err := s.GuildMemberRoleAdd(i.GuildID, userID, roleID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to add role.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Role added successfully!",
		},
	})
}

func removeRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var roleID, userID string

	for _, option := range i.ApplicationCommandData().Options[0].Options {
		switch option.Name {
		case "role":
			if option.Type != discordgo.ApplicationCommandOptionRole {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Invalid option type for role. Please provide a role.",
					},
				})
				return
			}
			roleID = option.RoleValue(s, "").ID
		case "user":
			if option.Type != discordgo.ApplicationCommandOptionUser {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Invalid option type for user. Please provide a user.",
					},
				})
				return
			}
			userID = option.UserValue(s).ID
		}
	}

	if roleID == "" || userID == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide both a role and a user.",
			},
		})
		return
	}

	err := s.GuildMemberRoleRemove(i.GuildID, userID, roleID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to remove role.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Role removed successfully!",
		},
	})
}

func roleAll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a role to add to everyone.",
			},
		})
		return
	}

	if options[0].Type != discordgo.ApplicationCommandOptionRole {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid option type. Please provide a role.",
			},
		})
		return
	}

	roleID := options[0].RoleValue(s, "").ID

	members, err := s.GuildMembers(i.GuildID, "", 1000)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve members.",
			},
		})
		return
	}

	var failedMembers []string
	for _, member := range members {
		err := s.GuildMemberRoleAdd(i.GuildID, member.User.ID, roleID)
		if err != nil {
			log.Printf("Failed to add role to member %s: %v", member.User.ID, err)
			failedMembers = append(failedMembers, member.User.ID)
			time.Sleep(1 * time.Second) // Adding a delay to handle potential rate limits
		}
	}

	if len(failedMembers) > 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to add role to some members: " + strings.Join(failedMembers, ", "),
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Role added to everyone successfully!",
			},
		})
	}
}
