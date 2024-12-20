package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CountingChannel struct {
	GuildID      string `bson:"guild_id"`
	ChannelID    string `bson:"channel_id"`
	LastCount    int    `bson:"last_count"`
	LastCountUID string `bson:"last_count_uid"`
}

func CountingCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if len(i.ApplicationCommandData().Options) == 0 {
		RespondWithMessage(s, i, "Please provide a subcommand.")
		return
	}

	switch i.ApplicationCommandData().Options[0].Name {
	case "set":
		setCountingChannel(s, i)
	case "get":
		getCountingChannel(s, i)
	case "delete":
		deleteCountingChannel(s, i)
	}
}

func setCountingChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel := i.ApplicationCommandData().Options[0].Options[0].ChannelValue(s)

	collection := db.GetCollection(cfg.DBName, "counting_channels")
	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"guild_id": i.GuildID},
		bson.M{"$set": bson.M{
			"channel_id": channel.ID,
			"last_count": 0,
		}},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		RespondWithMessage(s, i, "Failed to set counting channel.")
		return
	}

	RespondWithMessage(s, i, fmt.Sprintf("Counting channel set to <#%s>", channel.ID))
}

func getCountingChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var result CountingChannel
	err := db.GetCollection(cfg.DBName, "counting_channels").FindOne(
		context.TODO(),
		bson.M{"guild_id": i.GuildID},
	).Decode(&result)

	if err != nil {
		RespondWithMessage(s, i, "No counting channel set.")
		return
	}

	RespondWithMessage(s, i, fmt.Sprintf("Current counting channel: <#%s>", result.ChannelID))
}

func deleteCountingChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, err := db.GetCollection(cfg.DBName, "counting_channels").DeleteOne(
		context.TODO(),
		bson.M{"guild_id": i.GuildID},
	)

	if err != nil {
		RespondWithMessage(s, i, "Failed to delete counting channel.")
		return
	}

	RespondWithMessage(s, i, "Counting channel removed.")
}

func HandleCountingMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	var channel CountingChannel
	err := db.GetCollection(cfg.DBName, "counting_channels").FindOne(
		context.TODO(),
		bson.M{
			"guild_id":   m.GuildID,
			"channel_id": m.ChannelID,
		},
	).Decode(&channel)

	if err != nil {
		return
	}

	if channel.LastCountUID == m.Author.ID {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ %s, you can't count twice in a row!",
			m.Author.Mention()))
		return
	}

	number, err := strconv.Atoi(m.Content)
	if err != nil {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		return
	}

	if number != channel.LastCount+1 {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ %s tried to count %d, but we were at %d! The count has been reset to 1.",
			m.Author.Mention(), number, channel.LastCount))

		_, err = db.GetCollection(cfg.DBName, "counting_channels").UpdateOne(
			context.TODO(),
			bson.M{"guild_id": m.GuildID},
			bson.M{"$set": bson.M{
				"last_count":     0,
				"last_count_uid": "",
			}},
		)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to reset count.")
		}
		return
	}

	_, err = db.GetCollection(cfg.DBName, "counting_channels").UpdateOne(
		context.TODO(),
		bson.M{"guild_id": m.GuildID},
		bson.M{"$set": bson.M{
			"last_count":     number,
			"last_count_uid": m.Author.ID,
		}},
	)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to update count.")
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}
