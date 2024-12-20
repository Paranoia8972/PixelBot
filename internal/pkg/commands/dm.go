package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func DMCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		// Open modal for message input
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: fmt.Sprintf("dm_modal_%s", i.ApplicationCommandData().Options[0].UserValue(s).ID),
				Title:    "Send Direct Message",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:  "dm_message",
								Label:     "Message",
								Style:     discordgo.TextInputParagraph,
								Required:  true,
								MaxLength: 2000,
							},
						},
					},
				},
			},
		})
		if err != nil {
			RespondWithMessage(s, i, "Failed to create DM modal")
		}

	case discordgo.InteractionModalSubmit:
		// Handle modal submission
		data := i.ModalSubmitData()
		userID := strings.TrimPrefix(data.CustomID, "dm_modal_")
		message := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		// Create DM channel
		channel, err := s.UserChannelCreate(userID)
		if err != nil {
			RespondWithMessage(s, i, "Failed to create DM channel")
			return
		}

		// Send message
		_, err = s.ChannelMessageSend(channel.ID, message)
		if err != nil {
			RespondWithMessage(s, i, "Failed to send DM")
			return
		}

		RespondWithMessage(s, i, "Message sent successfully!")
	}
}
