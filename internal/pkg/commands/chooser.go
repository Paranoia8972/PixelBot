package commands

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ChooserCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options[0].Options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide either a list of items or a message ID.",
			},
		})
		return
	}

	subCmd := options[0]
	cmdOptions := subCmd.Options
	amount := 1

	for _, opt := range cmdOptions {
		if opt.Name == "amount" {
			amount = int(opt.IntValue())
			break
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if subCmd.Name == "message" {
		messageID := cmdOptions[0].StringValue()
		msg, err := s.ChannelMessage(i.ChannelID, messageID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Could not find message with that ID.",
				},
			})
			return
		}

		var users []string
		for _, reaction := range msg.Reactions {
			reactUsers, err := s.MessageReactions(i.ChannelID, messageID, reaction.Emoji.APIName(), 100, "", "")
			if err != nil {
				continue
			}
			for _, user := range reactUsers {
				if !user.Bot {
					users = append(users, user.ID)
				}
			}
		}

		userMap := make(map[string]bool)
		var uniqueUsers []string
		for _, user := range users {
			if !userMap[user] {
				userMap[user] = true
				uniqueUsers = append(uniqueUsers, user)
			}
		}

		if len(uniqueUsers) == 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "No users found who reacted to the message.",
				},
			})
			return
		}

		if amount > len(uniqueUsers) {
			amount = len(uniqueUsers)
		}

		chosen := make([]string, 0, amount)
		tempUsers := make([]string, len(uniqueUsers))
		copy(tempUsers, uniqueUsers)
		for i := 0; i < amount; i++ {
			idx := r.Intn(len(tempUsers))
			chosen = append(chosen, "<@"+tempUsers[idx]+">")
			tempUsers = append(tempUsers[:idx], tempUsers[idx+1:]...)
		}

		response := "Chosen user"
		if amount > 1 {
			response += "s"
		}
		response += " from reactions: " + strings.Join(chosen, ", ")

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
		return
	}

	items := strings.Split(cmdOptions[0].StringValue(), ",")
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}

	if len(items) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No valid items provided.",
			},
		})
		return
	}

	if amount > len(items) {
		amount = len(items)
	}

	chosen := make([]string, 0, amount)
	tempItems := make([]string, len(items))
	copy(tempItems, items)
	for i := 0; i < amount; i++ {
		idx := r.Intn(len(tempItems))
		chosen = append(chosen, tempItems[idx])
		tempItems = append(tempItems[:idx], tempItems[idx+1:]...)
	}

	response := "Chosen item"
	if amount > 1 {
		response += "s"
	}
	response += ": " + strings.Join(chosen, ", ")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
