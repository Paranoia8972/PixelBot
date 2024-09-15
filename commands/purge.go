package commands

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

func RemoveMessagesCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var count int
	var userID string

	for _, option := range options {
		switch option.Name {
		case "count":
			count = int(option.IntValue())
		case "user":
			userID = option.UserValue(nil).ID
		}
	}

	channelID := i.ChannelID
	messages, err := s.ChannelMessages(channelID, count, "", "", "")
	if err != nil {
		color.Red("Error fetching messages: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error fetching messages.",
				Flags:   64,
			},
		})
		return
	}

	var messageIDs []string
	twoWeeksAgo := time.Now().Add(-14 * 24 * time.Hour)
	var individuallyDeletedCount, bulkDeletedCount int

	for _, message := range messages {
		if userID == "" || message.Author.ID == userID {
			messageTime, err := time.Parse(time.RFC3339, message.Timestamp.Format(time.RFC3339))
			if err != nil {
				color.Red("Error parsing message timestamp: %v", err)
				continue
			}
			if messageTime.Before(twoWeeksAgo) {
				err = s.ChannelMessageDelete(channelID, message.ID)
				if err != nil {
					color.Red("Error deleting message: %v", err)
				} else {
					individuallyDeletedCount++
				}
			} else {
				messageIDs = append(messageIDs, message.ID)
			}
		}
	}

	if len(messageIDs) > 0 {
		err = s.ChannelMessagesBulkDelete(channelID, messageIDs)
		if err != nil {
			color.Red("Error deleting messages: %v", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error deleting messages.",
					Flags:   64,
				},
			})
			return
		}
		bulkDeletedCount = len(messageIDs)
	}

	totalDeleted := individuallyDeletedCount + bulkDeletedCount

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Successfully deleted " + strconv.Itoa(totalDeleted) + " messages.",
			Flags:   64,
		},
	})
}
