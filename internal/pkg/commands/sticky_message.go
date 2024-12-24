package commands

import (
	"context"
	"fmt"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StickyMessage struct {
	GuildID   string `bson:"guild_id"`
	ChannelID string `bson:"channel_id"`
	MessageID string `bson:"message_id"`
	Content   string `bson:"content"`
}

func StickyCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		if len(data.Options) > 0 {
			switch data.Options[0].Name {
			case "remove":
				RemoveStickyCommand(s, i)
				return
			case "list":
				ListStickyMessages(s, i)
				return
			}
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "sticky_modal",
				Title:    "Create Sticky Message",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:  "sticky_content",
								Label:     "Message Content",
								Style:     discordgo.TextInputParagraph,
								Required:  true,
								MaxLength: 2000,
							},
						},
					},
				},
			},
		})

	case discordgo.InteractionModalSubmit:
		data := i.ModalSubmitData()
		content := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		msg, err := s.ChannelMessageSend(i.ChannelID, content)
		if err != nil {
			RespondWithMessage(s, i, "Failed to send sticky message")
			return
		}

		collection := db.GetCollection(cfg.DBName, "sticky_messages")
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.M{"channel_id": i.ChannelID},
			bson.M{"$set": bson.M{
				"guild_id":   i.GuildID,
				"channel_id": i.ChannelID,
				"message_id": msg.ID,
				"content":    content,
			}},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			RespondWithMessage(s, i, "Failed to store sticky message")
			return
		}

		RespondWithMessage(s, i, "Sticky message created!")
	}
}

func HandleStickyMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	collection := db.GetCollection(cfg.DBName, "sticky_messages")
	var sticky StickyMessage
	err := collection.FindOne(
		context.TODO(),
		bson.M{"channel_id": m.ChannelID},
	).Decode(&sticky)

	if err != nil {
		return
	}

	s.ChannelMessageDelete(m.ChannelID, sticky.MessageID)

	msg, err := s.ChannelMessageSend(m.ChannelID, sticky.Content)
	if err != nil {
		return
	}

	collection.UpdateOne(
		context.TODO(),
		bson.M{"channel_id": m.ChannelID},
		bson.M{"$set": bson.M{"message_id": msg.ID}},
	)
}

func ListStickyMessages(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := i.GuildID
	filter := bson.M{"guild_id": guildID}
	cursor, err := db.GetCollection(cfg.DBName, "sticky_messages").Find(context.TODO(), filter)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve sticky messages.",
				Flags:   64,
			},
		})
		return
	}
	defer cursor.Close(context.TODO())

	var messages []StickyMessage
	if err = cursor.All(context.TODO(), &messages); err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to parse sticky messages.",
				Flags:   64,
			},
		})
		return
	}

	if len(messages) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No sticky messages found.",
				Flags:   64,
			},
		})
		return
	}

	var response string
	for _, msg := range messages {
		response += fmt.Sprintf("Channel: %s\nMessage: %s\n\n", msg.ChannelID, msg.Content)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   64,
		},
	})
}

func RemoveStickyCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		RespondWithMessage(s, i, "No message ID provided")
		return
	}

	subCommandOptions := options[0].Options
	if len(subCommandOptions) == 0 {
		RespondWithMessage(s, i, "No message ID provided")
		return
	}

	messageID := subCommandOptions[0].StringValue()

	s.ChannelMessageDelete(i.ChannelID, messageID)

	collection := db.GetCollection(cfg.DBName, "sticky_messages")
	result, err := collection.DeleteOne(
		context.TODO(),
		bson.M{"message_id": messageID},
	)
	if err != nil {
		RespondWithMessage(s, i, "Failed to remove sticky message")
		return
	}

	if result.DeletedCount == 0 {
		RespondWithMessage(s, i, "No sticky message found with that ID")
		return
	}

	RespondWithMessage(s, i, "Sticky message removed!")
}
