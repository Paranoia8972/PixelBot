package games

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
	LastMsgID    string `bson:"last_msg_id"`
	LastAuthorID string `bson:"last_author_id"`
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

	s.AddHandler(func(s *discordgo.Session, mrd *discordgo.MessageDeleteBulk) {
		for _, id := range mrd.Messages {
			if id == m.ID {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ %s deleted their number! Counting continues from %d.",
					m.Author.Mention(), channel.LastCount))
				return
			}
		}
	})

	s.AddHandler(func(s *discordgo.Session, md *discordgo.MessageDelete) {
		if md.ID == channel.LastMsgID && md.ID != m.ID {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Someone deleted the number %d! Counting continues from %d.",
				channel.LastCount, channel.LastCount))
		}
	})

	if channel.LastCountUID == m.Author.ID {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ %s, you can't count twice in a row! Wait for someone else to count.",
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
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ %s counted %d but the next number should have been %d! Counting restarted from 1.",
			m.Author.Mention(), number, channel.LastCount+1))

		_, err = db.GetCollection(cfg.DBName, "counting_channels").UpdateOne(
			context.TODO(),
			bson.M{"guild_id": m.GuildID},
			bson.M{"$set": bson.M{
				"last_count":     0,
				"last_count_uid": "",
				"last_msg_id":    "",
				"last_author_id": "",
			}},
		)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to reset count.")
		}
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	_, err = db.GetCollection(cfg.DBName, "counting_channels").UpdateOne(
		context.TODO(),
		bson.M{"guild_id": m.GuildID},
		bson.M{"$set": bson.M{
			"last_count":     number,
			"last_count_uid": m.Author.ID,
			"last_msg_id":    m.ID,
			"last_author_id": m.Author.ID,
		}},
	)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to update count.")
		return
	}

	if number%1000 == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("🎊 %s reached %d! What an achievement!",
			m.Author.Mention(), number))
	} else if number%500 == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("🎯 %s reached %d! Halfway to the next thousand!",
			m.Author.Mention(), number))
	} else if number%250 == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⭐ %s reached %d! Keep it up!",
			m.Author.Mention(), number))
	} else if number%100 == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("💯 %s reached %d! Nice work!",
			m.Author.Mention(), number))
		s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	}
}

func CountingDeleteHandler(s *discordgo.Session, md *discordgo.MessageDelete) {
	var channel CountingChannel
	err := db.GetCollection(cfg.DBName, "counting_channels").FindOne(
		context.TODO(),
		bson.M{"channel_id": md.ChannelID},
	).Decode(&channel)
	if err != nil {
		return
	}

	if md.ID == channel.LastMsgID {
		deletedNumber := channel.LastCount
		nextNumber := deletedNumber + 1

		s.ChannelMessageSend(md.ChannelID, fmt.Sprintf("❌ <@%s> deleted their number %d! The next number is %d.",
			channel.LastAuthorID, deletedNumber, nextNumber))
	}
}
