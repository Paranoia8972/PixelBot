package commands

import (
	"github.com/bwmarrin/discordgo"
)

func SayCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "say_modal",
				Title:    "Say Something",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "say_input",
								Label:       "Your Message",
								Style:       discordgo.TextInputParagraph,
								Placeholder: "Type your message here...",
								Required:    true,
							},
						},
					},
				},
			},
		})
	case discordgo.InteractionModalSubmit:

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Message sent!",
				Flags:   64,
			},
		})

		data := i.ModalSubmitData()
		message := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		s.ChannelMessageSend(i.ChannelID, message)
	}
}
