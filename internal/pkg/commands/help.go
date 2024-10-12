package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func HelpCommand(s *discordgo.Session, channelID string) {
	selectMenu := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				CustomID:    "select_menu",
				Placeholder: "Choose a command category",
				Options: []discordgo.SelectMenuOption{
					{
						Label:       "Ticket System",
						Value:       "ticket",
						Description: "Commands related to the ticket system",
					},
					{
						Label:       "Role Management",
						Value:       "role",
						Description: "Commands related to role management",
					},
					{
						Label:       "Giveaways",
						Value:       "giveaway",
						Description: "Commands related to giveaways",
					},
					{
						Label:       "General Commands",
						Value:       "general",
						Description: "General commands",
					},
				},
			},
		},
	}

	_, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: "Please select a command category:",
		Components: []discordgo.MessageComponent{
			selectMenu,
		},
	})
	if err != nil {
		log.Println("Error sending select menu:", err)
	}
}

func HandleSelectMenuInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	data := i.MessageComponentData()
	if data.CustomID != "select_menu" {
		return
	}

	var response string
	switch data.Values[0] {
	case "ticket":
		response = "Ticket System Commands:\n" +
			"- `/ticket setup`: Setup the ticket system\n" +
			"- `/ticket send`: Sends a new message with a button to create a ticket\n" +
			"- `/ticket close`: Closes an existing ticket"
	case "role":
		response = "Role Management Commands:\n" +
			"- `/role add`: Assign a role to a specific user\n" +
			"- `/role remove`: Remove a role from a specific user\n" +
			"- `/role all`: Assign a role to all members"
	case "giveaway":
		response = "Giveaway Commands:\n" +
			"- `/giveaway create`: Create a giveaway\n" +
			"- `/giveaway edit`: Edit a giveaway\n" +
			"- `/giveaway end`: End a giveaway\n" +
			"- `/giveaway reroll`: Reroll a giveaway"
	case "general":
		response = "General Commands:\n" +
			"- `/ping`: Responds with the Bot's latency\n" +
			"- `/say`: Repeats a message\n" +
			"- `/clear`: Deletes messages from a channel"
	default:
		response = "Unknown category"
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
	}
}
