package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func EditCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		messageID := i.ApplicationCommandData().Options[0].StringValue()

		msg, err := s.ChannelMessage(i.ChannelID, messageID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to retrieve the message.",
					Flags:   64,
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "edit_modal_" + messageID,
				Title:    "Edit Message",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID: "edit_input",
								Label:    "New Content",
								Style:    discordgo.TextInputParagraph,
								Value:    msg.Content,
								Required: true,
							},
						},
					},
				},
			},
		})

	case discordgo.InteractionModalSubmit:
		if !strings.HasPrefix(i.ModalSubmitData().CustomID, "edit_modal_") {
			return
		}
		messageID := strings.TrimPrefix(i.ModalSubmitData().CustomID, "edit_modal_")
		newContent := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		_, err := s.ChannelMessageEdit(i.ChannelID, messageID, newContent)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to edit the message.",
					Flags:   64,
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
	}
}
