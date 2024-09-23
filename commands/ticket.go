package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Paranoia8972/PixelBot/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/russross/blackfriday/v2"
)

func TicketCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	case "setup":
		TicketSetupCommand(s, i)
	case "addmodrole":
		TicketAddModeratorRolesCommand(s, i)
	case "getmodroles":
		TicketGetModeratorRolesCommand(s, i)
	case "removemodrole":
		TicketRemoveModeratorRolesCommand(s, i)
	default:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unknown subcommand.",
			},
		})
	}
}

func TicketSetupCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options
	if len(options) < 2 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a channel and category.",
			},
		})
		return
	}

	channelID := options[0].ChannelValue(s).ID
	categoryID := options[1].ChannelValue(s).ID
	transcriptChannelID := options[2].ChannelValue(s).ID

	err := utils.SetTicketSetup(i.GuildID, channelID, categoryID, transcriptChannelID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to set up ticket system.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ticket system set up successfully!",
		},
	})

	s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: "Click the button below to create a new ticket.",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "create_ticket",
						Label:    "Create Ticket",
						Style:    discordgo.PrimaryButton,
					},
				},
			},
		},
	})
}

func TicketButtonHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.MessageComponentData().CustomID != "create_ticket" {
		return
	}

	ticketChannel, err := utils.GetTicketSetup(i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve ticket channel.",
			},
		})
		return
	}

	username := i.Member.User.Username
	userID := i.Member.User.ID

	ticketNumber, err := utils.GetNextTicketNumber(i.GuildID, userID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve next ticket number.",
			},
		})
		return
	}

	channelName := "ticket-" + username + "-" + strconv.Itoa(ticketNumber)

	moderatorRoles, err := utils.GetModeratorRoles(i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve moderator roles.",
			},
		})
		return
	}

	permissionOverwrites := []*discordgo.PermissionOverwrite{
		// Set permissions for the ticket creator
		{
			ID:    userID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
		},
		// Set permissions for everyone else in the server
		{
			ID:   i.GuildID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
		},
	}

	// Set permissions for each moderator role
	for _, roleID := range moderatorRoles {
		permissionOverwrites = append(permissionOverwrites, &discordgo.PermissionOverwrite{
			ID:    roleID,
			Type:  discordgo.PermissionOverwriteTypeRole,
			Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
		})
	}

	channel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name:                 channelName,
		Type:                 discordgo.ChannelTypeGuildText,
		ParentID:             ticketChannel.CategoryID,
		PermissionOverwrites: permissionOverwrites,
	})
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create ticket channel.",
			},
		})
		return
	}

	err = utils.IncrementTicketNumber(i.GuildID, userID, ticketNumber)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to increment ticket number.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ticket created: <#" + channel.ID + ">",
		},
	})

	s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: "Ticket created by <@" + userID + ">",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "close_ticket",
						Label:    "Close Ticket",
						Style:    discordgo.DangerButton,
					},
				},
			},
		},
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Ticket Information",
				Description: "Ticket information for this ticket.",
				Color:       0x00ff00,
				Image: &discordgo.MessageEmbedImage{
					URL: "https://cdn.discordapp.com/attachments/1190806297881350164/1284438122867986492/footer.png",
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Ticket Creator",
						Value:  "<@" + userID + ">",
						Inline: true,
					},
					{
						Name:   "Ticket Number",
						Value:  strconv.Itoa(ticketNumber),
						Inline: true,
					},
				},
			},
		},
	})

}

func TicketAddModeratorRolesCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	err := utils.AddModeratorRoles(i.GuildID, []string{roleID})
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to add moderator role.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Moderator role added successfully!",
		},
	})
}

func TicketGetModeratorRolesCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	moderatorRoles, err := utils.GetModeratorRoles(i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve moderator roles.",
			},
		})
		return
	}

	var roleMentions []string
	for _, roleID := range moderatorRoles {
		roleMentions = append(roleMentions, "<@&"+roleID+">")
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Moderator Roles",
					Description: "Moderator roles for this server.",
					Color:       0x00ff00,
					Image: &discordgo.MessageEmbedImage{
						URL: "https://cdn.discordapp.com/attachments/1190806297881350164/1284438122867986492/footer.png",
					},
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

func TicketRemoveModeratorRolesCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	err := utils.RemoveModeratorRoles(i.GuildID, roleID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to remove moderator role.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Moderator role removed successfully!",
		},
	})
}

func TicketCloseHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		if i.MessageComponentData().CustomID == "close_ticket" {
			messages, err := s.ChannelMessages(i.ChannelID, 100, "", "", "")
			if err != nil {
				log.Printf("error fetching messages: %v", err)
				return
			}

			var transcript []map[string]interface{}
			for _, msg := range messages {
				messageData := map[string]interface{}{
					"username":        msg.Author.Username,
					"pfp":             msg.Author.AvatarURL(""),
					"message_content": string(blackfriday.Run([]byte(msg.Content))),
					"attachments":     []map[string]interface{}{},
					"embeds":          []map[string]interface{}{},
				}

				for _, attachment := range msg.Attachments {
					attachmentData := map[string]interface{}{
						"type":     attachment.ContentType,
						"url":      attachment.URL,
						"filename": attachment.Filename,
					}
					messageData["attachments"] = append(messageData["attachments"].([]map[string]interface{}), attachmentData)
				}

				for _, embed := range msg.Embeds {
					embedData := map[string]interface{}{
						"title":       "**" + embed.Title + "**",
						"description": embed.Description,
						"url":         embed.URL,
						"color":       embed.Color,
					}
					messageData["embeds"] = append(messageData["embeds"].([]map[string]interface{}), embedData)
				}

				transcript = append(transcript, messageData)
			}

			transcriptData := map[string]interface{}{
				"transcript": transcript,
			}

			transcriptJSON, err := json.Marshal(transcriptData)
			if err != nil {
				log.Printf("error marshalling transcript: %v", err)
				return
			}

			transcriptID, err := utils.StoreTranscript(i.GuildID, i.Member.User.ID, i.ChannelID, transcriptJSON)
			if err != nil {
				log.Printf("error storing transcript: %v", err)
				return
			}

			_, err = s.ChannelDelete(i.ChannelID)
			if err != nil {
				log.Printf("error deleting channel: %v", err)
			} else {
				log.Printf("channel %s deleted", i.ChannelID)
			}
			// DM the user the ticket ID
			channel, err := s.UserChannelCreate(i.Member.User.ID)
			if err != nil {
				log.Printf("error creating DM channel: %v", err)
				return
			}

			username := i.Member.User.Username
			userID := i.Member.User.ID

			ticketNumber, err := utils.GetNextTicketNumber(i.GuildID, userID)
			if err != nil {
				log.Printf("error getting next ticket number: %v", err)
				return
			}

			_, err = s.ChannelMessageSend(channel.ID, fmt.Sprintf("Your ticket "+"`ticket-"+username+"-"+strconv.Itoa(ticketNumber-1)+"` has been closed."+"\n\nHere is your transcript: https://ticket.ecyt.dev/ticket?id="+transcriptID.Hex()))
			if err != nil {
				log.Printf("error sending DM: %v", err)
				return
			}

			_, err = s.ChannelDelete(i.ChannelID)
			if err != nil {
				log.Printf("error deleting channel: %v", err)
			} else {
				log.Printf("channel %s deleted", i.ChannelID)
			}
		}
	}
}
