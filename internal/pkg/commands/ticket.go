package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/Paranoia8972/PixelBot/internal/pkg/utils"
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
	case "send":
		TicketSendMessage(s, i)
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
		respondWithMessage(s, i, "Please provide a channel and category.")
		return
	}

	channelID := options[0].ChannelValue(s).ID
	categoryID := options[1].ChannelValue(s).ID
	transcriptChannelID := options[2].ChannelValue(s).ID

	err := utils.SetTicketSetup(i.GuildID, channelID, categoryID, transcriptChannelID)
	if err != nil {
		respondWithMessage(s, i, "Failed to set up ticket system.")
		return
	}

	respondWithMessage(s, i, "Ticket system set up successfully!")

	s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: "Select an option below to open a ticket.",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "ticket_menu",
						Placeholder: "Choose an option",
						Options: []discordgo.SelectMenuOption{
							{
								Label: "Ban Appeal",
								Value: "ban_appeal",
								Emoji: &discordgo.ComponentEmoji{
									ID:   "1292525123596714048",
									Name: "ban",
								},
							},
							{
								Label: "Team Application",
								Value: "team_application",
								Emoji: &discordgo.ComponentEmoji{
									ID:   "1292526428147154994",
									Name: "team",
								},
							},
							{
								Label: "Bug Report",
								Value: "bug_report",
								Emoji: &discordgo.ComponentEmoji{
									ID:   "1292526495314608252",
									Name: "bug",
								},
							},
							{
								Label: "General support",
								Value: "general_support",
								Emoji: &discordgo.ComponentEmoji{
									ID:   "1292524894608822343",
									Name: "support",
								},
							},
						},
					},
				},
			},
		},
	})
}

func TicketSendMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {

	channelID, err := utils.GetChannelID(i.GuildID)
	if err != nil {
		respondWithMessage(s, i, "Error retrieving ChannelID.")
		return
	}

	respondWithMessage(s, i, "Ticket message sent successfully!")
	s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: "Select an option below to open a ticket.",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "ticket_menu",
						Placeholder: "Choose an option",
						Options: []discordgo.SelectMenuOption{
							{
								Label: "Ban Appeal",
								Value: "ban_appeal",
								Emoji: &discordgo.ComponentEmoji{
									ID:   "1292525123596714048",
									Name: "ban",
								},
							},
							{
								Label: "Team Application",
								Value: "team_application",
								Emoji: &discordgo.ComponentEmoji{
									ID:   "1292526428147154994",
									Name: "team",
								},
							},
							{
								Label: "Bug Report",
								Value: "bug_report",
								Emoji: &discordgo.ComponentEmoji{
									ID:   "1292526495314608252",
									Name: "bug",
								},
							},
							{
								Label: "General support",
								Value: "general_support",
								Emoji: &discordgo.ComponentEmoji{
									ID:   "1292524894608822343",
									Name: "support",
								},
							},
						},
					},
				},
			},
		},
	})
}

func TicketSelectHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	if data.CustomID != "ticket_menu" {
		log.Println("Error: CustomID does not match 'ticket_menu'")
		return
	}

	if len(data.Values) == 0 {
		log.Println("Error: No values selected in the select menu")
		return
	}

	selectedOption := data.Values[0]

	var modal *discordgo.InteractionResponse
	switch selectedOption {
	case "ban_appeal":
		modal = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "ban_appeal_modal",
				Title:    "Ban Appeal",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "ban_appeal_username",
								Label:       "Please provide your minecraft username",
								Style:       discordgo.TextInputShort,
								Placeholder: "Enter your username here...",
								Required:    true,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "ban_appeal_details",
								Label:       "Please provide details about your appeal",
								Style:       discordgo.TextInputParagraph,
								Placeholder: "Enter your appeal details here...",
								Required:    true,
							},
						},
					},
				},
			},
		}
	case "team_application":
		modal = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "team_application_modal",
				Title:    "Team Application",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "team_application_role",
								Label:       "Please provide the role you are applying for",
								Style:       discordgo.TextInputShort,
								Placeholder: "Enter the role here...",
								Required:    true,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "team_application_details",
								Label:       "Please provide details about your application",
								Style:       discordgo.TextInputParagraph,
								Placeholder: "Enter your application details here...",
								Required:    true,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "team_application_links",
								Label:       "Please provide any relevant links.",
								Style:       discordgo.TextInputShort,
								Placeholder: "Enter your links here... (github.com/...)",
								Required:    false,
							},
						},
					},
				},
			},
		}
	case "bug_report":
		modal = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "bug_report_modal",
				Title:    "Bug Report",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "bug_report_title",
								Label:       "Please provide a title for the bug report",
								Style:       discordgo.TextInputShort,
								Placeholder: "Enter the bug title here...",
								Required:    true,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "bug_report_details",
								Label:       "Please describe the bug you encountered",
								Style:       discordgo.TextInputParagraph,
								Placeholder: "Enter the bug details here...",
								Required:    true,
							},
						},
					},
				},
			},
		}
	case "general_support":
		modal = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "general_support_modal",
				Title:    "General Support",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "general_support_details",
								Label:       "How can we assist you?",
								Style:       discordgo.TextInputParagraph,
								Placeholder: "Enter your support request here...",
								Required:    true,
							},
						},
					},
				},
			},
		}
	}

	err := s.InteractionRespond(i.Interaction, modal)
	if err != nil {
		log.Println("Error responding to interaction:", err)
	}
}

func TicketModalSubmitHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionModalSubmit {
		log.Printf("ModalSubmitHandler called with incorrect interaction type: %v", i.Type)
		return
	}

	data := i.ModalSubmitData()
	var embed *discordgo.MessageEmbed

	switch data.CustomID {
	case "ban_appeal_modal":
		username := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		details := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		embed = &discordgo.MessageEmbed{
			Title:       "Ban Appeal Submitted",
			Description: "Your ban appeal has been submitted successfully.",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Username",
					Value: username,
				},
				{
					Name:  "Details",
					Value: details,
				},
			},
		}
	case "team_application_modal":
		username := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		details := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		links := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		embed = &discordgo.MessageEmbed{
			Title:       "Team Application Submitted",
			Description: "Your team application has been submitted successfully.",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Username",
					Value: username,
				},
				{
					Name:  "Details",
					Value: details,
				},
				{
					Name:  "Links",
					Value: links,
				},
			},
		}
	case "bug_report_modal":
		title := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		details := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		embed = &discordgo.MessageEmbed{
			Title:       "Bug Report Submitted",
			Description: "Your bug report has been submitted successfully.",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Title",
					Value: title,
				},
				{
					Name:  "Details",
					Value: details,
				},
			},
		}
	case "general_support_modal":
		details := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		embed = &discordgo.MessageEmbed{
			Title:       "Support Request Submitted",
			Description: "Your support request has been submitted successfully.",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Details",
					Value: details,
				},
			},
		}
	}

	ticketChannel, err := utils.GetTicketSetup(i.GuildID)
	if err != nil {
		respondWithMessage(s, i, "Failed to retrieve ticket channel.")
		return
	}

	username := i.Member.User.Username
	userID := i.Member.User.ID

	ticketNumber, err := utils.GetNextTicketNumber(i.GuildID, userID)
	if err != nil {
		respondWithMessage(s, i, "Failed to retrieve next ticket number.")
		return
	}

	channelName := "ticket-" + username + "-" + strconv.Itoa(ticketNumber)

	channel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name:     channelName,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: ticketChannel.CategoryID,
		Topic:    fmt.Sprintf("%s: %s", data.CustomID, embed.Fields[0].Value),
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    i.Member.User.ID, // user
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory,
				Deny:  0,
			},
			{
				ID:    i.GuildID, // everyone
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: 0,
				Deny:  discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory,
			},
			{
				ID:    "1074354995757076662", // replace with admin role id
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory,
				Deny:  0,
			},
		},
	})
	if err != nil {
		respondWithMessage(s, i, "Failed to create ticket channel.")
		return
	}

	details := make(map[string]string)
	for _, component := range data.Components {
		if actionRow, ok := component.(*discordgo.ActionsRow); ok {
			for _, item := range actionRow.Components {
				if input, ok := item.(*discordgo.TextInput); ok {
					details[input.CustomID] = input.Value
				}
			}
		}
	}

	ticketType := data.CustomID

	ticket := utils.Tickets{
		GuildID:   i.GuildID,
		UserID:    userID,
		Username:  username,
		ChannelID: channel.ID,
		Type:      ticketType,
		Details:   details,
	}
	_, err = utils.StoreTicket(ticket)
	if err != nil {
		respondWithMessage(s, i, "Failed to store ticket information.")
		return
	}

	err = utils.IncrementTicketNumber(i.GuildID, userID, ticketNumber)
	if err != nil {
		respondWithMessage(s, i, "Failed to increment ticket number.")
		return
	}

	respondWithMessage(s, i, "Ticket created: <#"+channel.ID+">")
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
		Embeds: []*discordgo.MessageEmbed{embed},
	})
}

func createMessageData(msg *discordgo.Message) map[string]interface{} {
	messageData := map[string]interface{}{
		"username":        msg.Author.Username,
		"pfp":             msg.Author.AvatarURL(""),
		"message_content": string(blackfriday.Run([]byte(msg.Content))),
		"attachments":     []map[string]interface{}{},
		"embeds":          []map[string]interface{}{},
		"reactions":       []map[string]interface{}{},
		"components":      []map[string]interface{}{},
	}

	for _, reaction := range msg.Reactions {
		reactionData := map[string]interface{}{
			"emoji": reaction.Emoji.Name,
			"count": reaction.Count,
		}
		messageData["reactions"] = append(messageData["reactions"].([]map[string]interface{}), reactionData)
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
			"title":       embed.Title,
			"description": embed.Description,
			"url":         embed.URL,
			"color":       embed.Color,
			"fields":      []map[string]interface{}{},
			"image":       embed.Image,
		}
		for _, field := range embed.Fields {
			fieldData := map[string]interface{}{
				"name":  field.Name,
				"value": field.Value,
			}
			embedData["fields"] = append(embedData["fields"].([]map[string]interface{}), fieldData)
		}
		messageData["embeds"] = append(messageData["embeds"].([]map[string]interface{}), embedData)
	}

	return messageData
}

func TicketCloseHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "close_ticket" {
		messages, err := s.ChannelMessages(i.ChannelID, 100, "", "", "")
		if err != nil {
			log.Printf("error fetching messages: %v", err)
			return
		}

		var transcript []map[string]interface{}
		for _, msg := range messages {
			messageData := createMessageData(msg)
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

		if _, err := s.ChannelDelete(i.ChannelID); err != nil {
			log.Printf("error deleting channel: %v", err)
		} else {
			log.Printf("channel %s deleted", i.ChannelID)
		}

		ticket, err := utils.GetTicketByChannelID(i.ChannelID)
		if err != nil {
			log.Printf("error retrieving ticket: %v", err)
			return
		}

		channel, err := s.UserChannelCreate(ticket.UserID)
		if err != nil {
			log.Printf("error creating DM channel: %v", err)
			return
		}

		username := ticket.Username
		userID := ticket.UserID

		ticketNumber, err := utils.GetNextTicketNumber(i.GuildID, userID)
		if err != nil {
			log.Printf("error getting next ticket number: %v", err)
			return
		}

		message := fmt.Sprintf("Your ticket `ticket-%s-%d` has been closed.\n\nHere is your transcript: https://%s/ticket?id=%s",
			username, ticketNumber-1, cfg.TranscriptUrl, transcriptID.Hex())

		if _, err := s.ChannelMessageSend(channel.ID, message); err != nil {
			log.Printf("error sending DM: %v", err)
		}
		transcriptChannelID, err := utils.GetTranscriptChannelID(i.GuildID)
		if err != nil {
			log.Printf("error retrieving transcript channel ID: %v", err)
			return
		}

		details := ""
		for key, value := range ticket.Details {
			details += fmt.Sprintf("%s: %v\n", key, value)
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Ticket Closed",
			Description: fmt.Sprintf("Ticket Type: %s\nTranscript URL: [Click Here](https://%s/ticket?id=%s)\nTicket Channel: %s\nDetails:", ticket.Type, cfg.TranscriptUrl, transcriptID.Hex(), i.ChannelID),
			Color:       0xFF0000,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Details",
					Value: details,
				},
			},
		}

		if _, err := s.ChannelMessageSendEmbed(transcriptChannelID, embed); err != nil {
			log.Printf("error sending embed to transcript channel: %v", err)
		}

	}
}
